# Markdown Parser Test

This is a comprehensive test of the **markdown parser** written in _Go_.

## Features

The parser supports the following elements:

- **Headings** (h1-h6)
- _Paragraphs_ with inline formatting
- Unordered lists
- Ordered lists
- Links like [Google](https://www.google.com)
- Inline `code` snippets
- Code blocks

### Code Example

Here's a simple Go function:

```go
func hello() {
    fmt.Println("Hello, World!")
}
```

## Lists

### Unordered List

- First item
- Second item with **bold text**
- Third item with _italic text_

### Ordered List

1. First ordered item
2. Second ordered item with `inline code`
3. Third ordered item

## Links and Formatting

Visit [GitHub](https://github.com) for more Go projects.

You can combine **bold** and _italic_ formatting, and even include `inline code` within paragraphs.

## Conclusion

This markdown parser demonstrates basic markdown-to-HTML conversion capabilities in Go.
