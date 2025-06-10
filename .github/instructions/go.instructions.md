---
applyTo: "**/*.go"
---
# Go Coding Guidelines

- All exported identifiers (types, functions, methods, variables, constants) must have clear and concise English comments.
- For interfaces, provide detailed comments for each method in the interface definition, describing its purpose and usage, especially in the context of the domain (e.g., graphics).
- For interface implementations, only use a brief comment like `// <function> implements <InterfaceName>.`
- Do not redefine or shadow Go built-in functions such as `max` and `min`.
- Follow Go best practices for readability, maintainability, and simplicity.
- Use standard libraries and idiomatic Go patterns.
- Keep the order of implementation methods consistent with the order in the interface definition.
