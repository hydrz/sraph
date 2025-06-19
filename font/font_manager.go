// Package font - Font manager functionality
package font

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"sync"
)

// FontManager manages font registration, discovery, and caching.
// It provides a centralized way to access fonts throughout the application.
type FontManager interface {
	// RegisterFont registers a font with the manager.
	RegisterFont(font Font, name string) error

	// RegisterFontFromFile loads and registers a font from a file.
	RegisterFontFromFile(path string, name string) error

	// GetFont retrieves a font by name.
	GetFont(name string) (Font, bool)

	// GetFontByFamily retrieves a font by family, weight, and style.
	GetFontByFamily(family string, weight Weight, style Style) (Font, bool)

	// GetRegisteredFonts returns all registered font names.
	GetRegisteredFonts() []string

	// GetFontFamilies returns all available font families.
	GetFontFamilies() []string

	// SetDefaultFont sets the default font.
	SetDefaultFont(font Font)

	// GetDefaultFont returns the default font.
	GetDefaultFont() Font

	// GetFallbackFonts returns fallback fonts for the given font.
	GetFallbackFonts(font Font) []Font

	// ClearCache clears the font cache.
	ClearCache()
}

// FontDescriptor describes a font for loading and matching.
type FontDescriptor struct {
	// FamilyName is the font family name.
	FamilyName string
	// Weight is the font weight.
	Weight Weight
	// Style is the font style.
	Style Style
	// Size is the font size in points.
	Size float32
	// AxisAlignment is the subpixel alignment.
	AxisAlignment AxisAlignment
}

// FontCache provides font caching functionality.
type FontCache interface {
	// Get retrieves a font from the cache.
	Get(desc FontDescriptor) (Font, bool)

	// Put stores a font in the cache.
	Put(desc FontDescriptor, font Font)

	// Clear clears the cache.
	Clear()

	// Size returns the number of cached fonts.
	Size() int
}

// defaultFontManager provides the default implementation of FontManager.
type defaultFontManager struct {
	fonts       map[string]Font
	families    map[string]map[Weight]map[Style]Font
	defaultFont Font
	fallbacks   map[uint64][]Font
	cache       FontCache
	loader      TypefaceLoader
	mutex       sync.RWMutex
}

// NewFontManager creates a new font manager.
func NewFontManager() FontManager {
	return &defaultFontManager{
		fonts:     make(map[string]Font),
		families:  make(map[string]map[Weight]map[Style]Font),
		fallbacks: make(map[uint64][]Font),
		cache:     NewFontCache(),
		loader:    NewTypefaceLoader(),
	}
}

// RegisterFont implements FontManager.
func (m *defaultFontManager) RegisterFont(font Font, name string) error {
	if font == nil || !font.IsValid() {
		return errors.New("invalid font")
	}

	if name == "" {
		name = fmt.Sprintf("font_%d", font.GetHash())
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.fonts[name] = font

	// Register in families if possible
	if typeface := font.GetTypeface(); typeface != nil {
		if dt, ok := typeface.(*defaultTypeface); ok {
			info := dt.GetInfo()
			if m.families[info.FamilyName] == nil {
				m.families[info.FamilyName] = make(map[Weight]map[Style]Font)
			}
			if m.families[info.FamilyName][info.Weight] == nil {
				m.families[info.FamilyName][info.Weight] = make(map[Style]Font)
			}
			m.families[info.FamilyName][info.Weight][info.Style] = font
		}
	}

	return nil
}

// RegisterFontFromFile implements FontManager.
func (m *defaultFontManager) RegisterFontFromFile(path string, name string) error {
	typeface, err := m.loader.LoadFromFile(path)
	if err != nil {
		return fmt.Errorf("failed to load font from file %s: %w", path, err)
	}

	// Extract size from filename or use default
	size := extractSizeFromPath(path)
	if size <= 0 {
		size = 16 // Default size
	}

	font := DefaultFont(typeface, size, AxisAlignmentNone)
	if font == nil {
		return errors.New("failed to create font")
	}

	if name == "" {
		name = strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	}

	return m.RegisterFont(font, name)
}

// GetFont implements FontManager.
func (m *defaultFontManager) GetFont(name string) (Font, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	font, exists := m.fonts[name]
	return font, exists
}

// GetFontByFamily implements FontManager.
func (m *defaultFontManager) GetFontByFamily(family string, weight Weight, style Style) (Font, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if families, exists := m.families[family]; exists {
		if weights, exists := families[weight]; exists {
			if font, exists := weights[style]; exists {
				return font, true
			}
		}

		// Fallback to closest weight
		for w := range families {
			if styles, exists := families[w]; exists {
				if font, exists := styles[style]; exists {
					return font, true
				}
			}
		}

		// Fallback to any style
		for _, weights := range families {
			for _, font := range weights {
				return font, true
			}
		}
	}

	return nil, false
}

// GetRegisteredFonts implements FontManager.
func (m *defaultFontManager) GetRegisteredFonts() []string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	names := make([]string, 0, len(m.fonts))
	for name := range m.fonts {
		names = append(names, name)
	}
	return names
}

// GetFontFamilies implements FontManager.
func (m *defaultFontManager) GetFontFamilies() []string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	families := make([]string, 0, len(m.families))
	for family := range m.families {
		families = append(families, family)
	}
	return families
}

// SetDefaultFont implements FontManager.
func (m *defaultFontManager) SetDefaultFont(font Font) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.defaultFont = font
}

// GetDefaultFont implements FontManager.
func (m *defaultFontManager) GetDefaultFont() Font {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.defaultFont
}

// GetFallbackFonts implements FontManager.
func (m *defaultFontManager) GetFallbackFonts(font Font) []Font {
	if font == nil {
		return nil
	}

	m.mutex.RLock()
	defer m.mutex.RUnlock()

	fallbacks, exists := m.fallbacks[font.GetHash()]
	if exists {
		return fallbacks
	}

	return nil
}

// ClearCache implements FontManager.
func (m *defaultFontManager) ClearCache() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.cache.Clear()
}

// defaultFontCache provides the default implementation of FontCache.
type defaultFontCache struct {
	fonts map[string]Font
	mutex sync.RWMutex
}

// NewFontCache creates a new font cache.
func NewFontCache() FontCache {
	return &defaultFontCache{
		fonts: make(map[string]Font),
	}
}

// Get implements FontCache.
func (c *defaultFontCache) Get(desc FontDescriptor) (Font, bool) {
	key := c.makeKey(desc)
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	font, exists := c.fonts[key]
	return font, exists
}

// Put implements FontCache.
func (c *defaultFontCache) Put(desc FontDescriptor, font Font) {
	key := c.makeKey(desc)
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.fonts[key] = font
}

// Clear implements FontCache.
func (c *defaultFontCache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.fonts = make(map[string]Font)
}

// Size implements FontCache.
func (c *defaultFontCache) Size() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return len(c.fonts)
}

// makeKey creates a cache key from a font descriptor.
func (c *defaultFontCache) makeKey(desc FontDescriptor) string {
	return fmt.Sprintf("%s_%d_%d_%.2f_%d", desc.FamilyName, desc.Weight, desc.Style, desc.Size, desc.AxisAlignment)
}

// Helper functions

func extractSizeFromPath(path string) float32 {
	// Try to extract size from filename
	// This is a simple implementation, could be more sophisticated
	name := strings.ToLower(filepath.Base(path))

	// Common size patterns
	if strings.Contains(name, "12pt") || strings.Contains(name, "12") {
		return 12
	}
	if strings.Contains(name, "14pt") || strings.Contains(name, "14") {
		return 14
	}
	if strings.Contains(name, "16pt") || strings.Contains(name, "16") {
		return 16
	}
	if strings.Contains(name, "18pt") || strings.Contains(name, "18") {
		return 18
	}
	if strings.Contains(name, "24pt") || strings.Contains(name, "24") {
		return 24
	}

	return 0
}
