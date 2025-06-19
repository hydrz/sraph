// Package font - Typeface loading functionality
package font

import (
	"errors"
	"fmt"
	"hash/fnv"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/opensraph/sraph/geom"
)

// TypefaceLoader handles loading typefaces from various sources.
type TypefaceLoader interface {
	// LoadFromFile loads a typeface from a font file.
	LoadFromFile(path string) (Typeface, error)

	// LoadFromData loads a typeface from font data in memory.
	LoadFromData(data []byte) (Typeface, error)

	// GetRegisteredTypefaces returns all registered typefaces.
	GetRegisteredTypefaces() []Typeface

	// RegisterTypeface registers a typeface with an optional name.
	RegisterTypeface(typeface Typeface, name string) error

	// FindTypeface finds a typeface by family name and style.
	FindTypeface(familyName string, weight Weight, style Style) (Typeface, bool)
}

// FontFormat represents supported font file formats.
type FontFormat int

const (
	// FontFormatTTF represents TrueType Font format.
	FontFormatTTF FontFormat = iota
	// FontFormatOTF represents OpenType Font format.
	FontFormatOTF
	// FontFormatWOFF represents Web Open Font Format.
	FontFormatWOFF
	// FontFormatWOFF2 represents Web Open Font Format 2.
	FontFormatWOFF2
)

// TypefaceInfo contains metadata about a loaded typeface.
type TypefaceInfo struct {
	// FamilyName is the font family name.
	FamilyName string
	// StyleName is the font style name.
	StyleName string
	// Weight is the font weight.
	Weight Weight
	// Style is the font style.
	Style Style
	// Format is the font file format.
	Format FontFormat
	// FilePath is the original file path (if loaded from file).
	FilePath string
	// IsValid indicates if the typeface is valid.
	IsValid bool
}

// defaultTypeface provides a basic implementation of Typeface.
type defaultTypeface struct {
	info         TypefaceInfo
	data         []byte
	hash         uint64
	glyphCount   int
	unitsPerEM   int
	ascender     float32
	descender    float32
	capHeight    float32
	xHeight      float32
	glyphCache   map[rune]GlyphIndex
	metricsCache map[GlyphIndex]GlyphMetrics
	mutex        sync.RWMutex
}

// NewDefaultTypeface creates a new default typeface implementation.
func NewDefaultTypeface(data []byte, info TypefaceInfo) Typeface {
	hash := calculateDataHash(data)

	return &defaultTypeface{
		info:         info,
		data:         data,
		hash:         hash,
		glyphCache:   make(map[rune]GlyphIndex),
		metricsCache: make(map[GlyphIndex]GlyphMetrics),
	}
}

// IsValid implements Typeface.
func (t *defaultTypeface) IsValid() bool {
	return t.info.IsValid && len(t.data) > 0
}

// GetHash implements Typeface.
func (t *defaultTypeface) GetHash() uint64 {
	return t.hash
}

// GetGlyphCount implements Typeface.
func (t *defaultTypeface) GetGlyphCount() int {
	return t.glyphCount
}

// GetUnitsPerEM implements Typeface.
func (t *defaultTypeface) GetUnitsPerEM() int {
	return t.unitsPerEM
}

// GetAscender implements Typeface.
func (t *defaultTypeface) GetAscender() float32 {
	return t.ascender
}

// GetDescender implements Typeface.
func (t *defaultTypeface) GetDescender() float32 {
	return t.descender
}

// GetCapHeight implements Typeface.
func (t *defaultTypeface) GetCapHeight() float32 {
	return t.capHeight
}

// GetXHeight implements Typeface.
func (t *defaultTypeface) GetXHeight() float32 {
	return t.xHeight
}

// GetGlyphForCodepoint implements Typeface.
func (t *defaultTypeface) GetGlyphForCodepoint(codepoint rune) (GlyphIndex, bool) {
	t.mutex.RLock()
	if glyph, exists := t.glyphCache[codepoint]; exists {
		t.mutex.RUnlock()
		return glyph, true
	}
	t.mutex.RUnlock()

	// TODO: Implement actual glyph lookup from font data
	// For now, return a placeholder
	glyph := GlyphIndex(uint16(codepoint) % 1000)

	t.mutex.Lock()
	t.glyphCache[codepoint] = glyph
	t.mutex.Unlock()

	return glyph, true
}

// GetGlyphMetrics implements Typeface.
func (t *defaultTypeface) GetGlyphMetrics(glyph GlyphIndex) GlyphMetrics {
	t.mutex.RLock()
	if metrics, exists := t.metricsCache[glyph]; exists {
		t.mutex.RUnlock()
		return metrics
	}
	t.mutex.RUnlock()

	// TODO: Implement actual metrics calculation from font data
	// For now, return placeholder metrics
	metrics := GlyphMetrics{
		AdvanceWidth:     float32(t.unitsPerEM) * 0.6,
		AdvanceHeight:    float32(t.unitsPerEM),
		LeftSideBearing:  0,
		RightSideBearing: 0,
		BoundingBox:      geom.Rect[geom.F32]{Left: geom.F32(0), Top: geom.F32(0), Right: geom.F32(float32(t.unitsPerEM) * 0.6), Bottom: geom.F32(t.unitsPerEM)},
	}

	t.mutex.Lock()
	t.metricsCache[glyph] = metrics
	t.mutex.Unlock()

	return metrics
}

// GetKerning implements Typeface.
func (t *defaultTypeface) GetKerning(left, right GlyphIndex) float32 {
	// TODO: Implement actual kerning lookup from font data
	return 0
}

// GetInfo returns the typeface information.
func (t *defaultTypeface) GetInfo() TypefaceInfo {
	return t.info
}

// defaultTypefaceLoader provides the default implementation of TypefaceLoader.
type defaultTypefaceLoader struct {
	typefaces map[string]Typeface
	families  map[string]map[Weight]map[Style]Typeface
	mutex     sync.RWMutex
}

// NewTypefaceLoader creates a new typeface loader.
func NewTypefaceLoader() TypefaceLoader {
	return &defaultTypefaceLoader{
		typefaces: make(map[string]Typeface),
		families:  make(map[string]map[Weight]map[Style]Typeface),
	}
}

// LoadFromFile implements TypefaceLoader.
func (l *defaultTypefaceLoader) LoadFromFile(path string) (Typeface, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read font file %s: %w", path, err)
	}

	format := detectFontFormat(path, data)
	info := TypefaceInfo{
		FamilyName: extractFamilyName(path),
		StyleName:  extractStyleName(path),
		Weight:     extractWeight(path),
		Style:      extractStyle(path),
		Format:     format,
		FilePath:   path,
		IsValid:    true,
	}

	return NewDefaultTypeface(data, info), nil
}

// LoadFromData implements TypefaceLoader.
func (l *defaultTypefaceLoader) LoadFromData(data []byte) (Typeface, error) {
	if len(data) == 0 {
		return nil, errors.New("empty font data")
	}

	// TODO: Parse font data to extract metadata
	info := TypefaceInfo{
		FamilyName: "Unknown",
		StyleName:  "Regular",
		Weight:     WeightNormal,
		Style:      StyleNormal,
		Format:     FontFormatTTF,
		IsValid:    true,
	}

	typeface := NewDefaultTypeface(data, info)
	return typeface, nil
}

// GetRegisteredTypefaces implements TypefaceLoader.
func (l *defaultTypefaceLoader) GetRegisteredTypefaces() []Typeface {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	typefaces := make([]Typeface, 0, len(l.typefaces))
	for _, typeface := range l.typefaces {
		typefaces = append(typefaces, typeface)
	}
	return typefaces
}

// RegisterTypeface implements TypefaceLoader.
func (l *defaultTypefaceLoader) RegisterTypeface(typeface Typeface, name string) error {
	if typeface == nil || !typeface.IsValid() {
		return errors.New("invalid typeface")
	}

	if name == "" {
		name = fmt.Sprintf("typeface_%d", typeface.GetHash())
	}

	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.typefaces[name] = typeface

	// Register in families map if it's a defaultTypeface
	if dt, ok := typeface.(*defaultTypeface); ok {
		info := dt.GetInfo()
		if l.families[info.FamilyName] == nil {
			l.families[info.FamilyName] = make(map[Weight]map[Style]Typeface)
		}
		if l.families[info.FamilyName][info.Weight] == nil {
			l.families[info.FamilyName][info.Weight] = make(map[Style]Typeface)
		}
		l.families[info.FamilyName][info.Weight][info.Style] = typeface
	}

	return nil
}

// FindTypeface implements TypefaceLoader.
func (l *defaultTypefaceLoader) FindTypeface(familyName string, weight Weight, style Style) (Typeface, bool) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	if family, exists := l.families[familyName]; exists {
		if weights, exists := family[weight]; exists {
			if typeface, exists := weights[style]; exists {
				return typeface, true
			}
		}

		// Fallback to closest weight
		for w := range family {
			if styles, exists := family[w]; exists {
				if typeface, exists := styles[style]; exists {
					return typeface, true
				}
			}
		}

		// Fallback to any style
		for _, weights := range family {
			for _, typeface := range weights {
				return typeface, true
			}
		}
	}

	return nil, false
}

// Helper functions

func calculateDataHash(data []byte) uint64 {
	h := fnv.New64a()
	h.Write(data)
	return h.Sum64()
}

func detectFontFormat(path string, data []byte) FontFormat {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".ttf":
		return FontFormatTTF
	case ".otf":
		return FontFormatOTF
	case ".woff":
		return FontFormatWOFF
	case ".woff2":
		return FontFormatWOFF2
	default:
		// Try to detect from data
		if len(data) >= 4 {
			magic := string(data[:4])
			switch magic {
			case "OTTO":
				return FontFormatOTF
			case "wOFF":
				return FontFormatWOFF
			case "wOF2":
				return FontFormatWOFF2
			default:
				return FontFormatTTF
			}
		}
		return FontFormatTTF
	}
}

func extractFamilyName(path string) string {
	name := filepath.Base(path)
	name = strings.TrimSuffix(name, filepath.Ext(name))

	// Remove common style suffixes
	styleSuffixes := []string{"-Regular", "-Bold", "-Italic", "-BoldItalic", "_Regular", "_Bold", "_Italic", "_BoldItalic"}
	for _, suffix := range styleSuffixes {
		name = strings.TrimSuffix(name, suffix)
	}

	return name
}

func extractStyleName(path string) string {
	name := filepath.Base(path)
	name = strings.TrimSuffix(name, filepath.Ext(name))

	if strings.Contains(strings.ToLower(name), "bolditalic") {
		return "Bold Italic"
	} else if strings.Contains(strings.ToLower(name), "bold") {
		return "Bold"
	} else if strings.Contains(strings.ToLower(name), "italic") {
		return "Italic"
	}

	return "Regular"
}

func extractWeight(path string) Weight {
	name := strings.ToLower(filepath.Base(path))

	if strings.Contains(name, "thin") {
		return WeightThin
	} else if strings.Contains(name, "extralight") || strings.Contains(name, "ultralight") {
		return WeightExtraLight
	} else if strings.Contains(name, "light") {
		return WeightLight
	} else if strings.Contains(name, "medium") {
		return WeightMedium
	} else if strings.Contains(name, "semibold") || strings.Contains(name, "demibold") {
		return WeightSemiBold
	} else if strings.Contains(name, "extrabold") || strings.Contains(name, "ultrabold") {
		return WeightExtraBold
	} else if strings.Contains(name, "bold") {
		return WeightBold
	} else if strings.Contains(name, "black") || strings.Contains(name, "heavy") {
		return WeightBlack
	}

	return WeightNormal
}

func extractStyle(path string) Style {
	name := strings.ToLower(filepath.Base(path))

	if strings.Contains(name, "italic") {
		return StyleItalic
	} else if strings.Contains(name, "oblique") {
		return StyleOblique
	}

	return StyleNormal
}
