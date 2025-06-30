package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strings"
)

// MarkdownElement represents different types of markdown elements
type MarkdownElement interface {
	ToHTML() string
}

// Heading represents markdown headings (h1-h6)
type Heading struct {
	Level int
	Text  string
}

// Paragraph represents markdown paragraphs
type Paragraph struct {
	Text string
}

// List represents markdown lists (ordered and unordered)
type List struct {
	Ordered bool
	Items   []string
}

// Link represents markdown links
type Link struct {
	Text string
	URL  string
}

// CodeBlock represents markdown code blocks
type CodeBlock struct {
	Language string
	Code     string
}

// InlineCode represents inline code
type InlineCode struct {
	Code string
}

// MarkdownParser represents the parser
type MarkdownParser struct {
	lines []string
}

// NewMarkdownParser creates a new parser instance
func NewMarkdownParser(lines []string) *MarkdownParser {
	return &MarkdownParser{lines: lines}
}

// Parse parses the markdown content and returns HTML elements
func (p *MarkdownParser) Parse() []MarkdownElement {
	var elements []MarkdownElement
	i := 0

	for i < len(p.lines) {
		line := strings.TrimSpace(p.lines[i])

		// Skip empty lines
		if line == "" {
			i++
			continue
		}

		// Parse headings
		if strings.HasPrefix(line, "#") {
			element, consumed := p.parseHeading(i)
			if element != nil {
				elements = append(elements, element)
			}
			i += consumed
			continue
		}

		// Parse code blocks
		if strings.HasPrefix(line, "```") {
			element, consumed := p.parseCodeBlock(i)
			if element != nil {
				elements = append(elements, element)
			}
			i += consumed
			continue
		}

		// Parse lists
		if strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "* ") ||
			regexp.MustCompile(`^\d+\. `).MatchString(line) {
			element, consumed := p.parseList(i)
			if element != nil {
				elements = append(elements, element)
			}
			i += consumed
			continue
		}

		// Parse paragraphs (default case)
		element, consumed := p.parseParagraph(i)
		if element != nil {
			elements = append(elements, element)
		}
		i += consumed
	}

	return elements
}

// parseHeading parses markdown headings
func (p *MarkdownParser) parseHeading(start int) (MarkdownElement, int) {
	line := strings.TrimSpace(p.lines[start])
	level := 0

	for i, char := range line {
		if char == '#' {
			level++
		} else {
			text := strings.TrimSpace(line[i:])
			return Heading{Level: level, Text: p.processInlineElements(text)}, 1
		}
	}

	return nil, 1
}

// parseCodeBlock parses markdown code blocks
func (p *MarkdownParser) parseCodeBlock(start int) (MarkdownElement, int) {
	firstLine := strings.TrimSpace(p.lines[start])
	language := strings.TrimSpace(firstLine[3:]) // Remove ```

	var codeLines []string
	i := start + 1

	for i < len(p.lines) {
		line := p.lines[i]
		if strings.TrimSpace(line) == "```" {
			break
		}
		codeLines = append(codeLines, line)
		i++
	}

	code := strings.Join(codeLines, "\n")
	return CodeBlock{Language: language, Code: code}, i - start + 1
}

// parseList parses markdown lists
func (p *MarkdownParser) parseList(start int) (MarkdownElement, int) {
	var items []string
	i := start
	isOrdered := regexp.MustCompile(`^\d+\. `).MatchString(strings.TrimSpace(p.lines[start]))

	for i < len(p.lines) {
		line := strings.TrimSpace(p.lines[i])

		if line == "" {
			i++
			continue
		}

		// Check if this is a list item
		if strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "* ") {
			if isOrdered {
				break // Different list type
			}
			items = append(items, p.processInlineElements(line[2:]))
		} else if regexp.MustCompile(`^\d+\. `).MatchString(line) {
			if !isOrdered {
				break // Different list type
			}
			re := regexp.MustCompile(`^\d+\. (.*)`)
			matches := re.FindStringSubmatch(line)
			if len(matches) > 1 {
				items = append(items, p.processInlineElements(matches[1]))
			}
		} else {
			break // Not a list item
		}
		i++
	}

	return List{Ordered: isOrdered, Items: items}, i - start
}

// parseParagraph parses markdown paragraphs
func (p *MarkdownParser) parseParagraph(start int) (MarkdownElement, int) {
	var paragraphLines []string
	i := start

	for i < len(p.lines) {
		line := strings.TrimSpace(p.lines[i])

		if line == "" {
			break
		}

		// Stop if we encounter other markdown elements
		if strings.HasPrefix(line, "#") ||
			strings.HasPrefix(line, "```") ||
			strings.HasPrefix(line, "- ") ||
			strings.HasPrefix(line, "* ") ||
			regexp.MustCompile(`^\d+\. `).MatchString(line) {
			break
		}

		paragraphLines = append(paragraphLines, line)
		i++
	}

	if len(paragraphLines) == 0 {
		return nil, 1
	}

	text := strings.Join(paragraphLines, " ")
	return Paragraph{Text: p.processInlineElements(text)}, i - start
}

// processInlineElements processes inline markdown elements like bold, italic, links, and inline code
func (p *MarkdownParser) processInlineElements(text string) string {
	// Process inline code first (to avoid conflicts with other formatting)
	re := regexp.MustCompile("`([^`]+)`")
	text = re.ReplaceAllString(text, "<code>$1</code>")

	// Process links [text](url)
	re = regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)
	text = re.ReplaceAllString(text, "<a href=\"$2\">$1</a>")

	// Process bold **text**
	re = regexp.MustCompile(`\*\*([^*]+)\*\*`)
	text = re.ReplaceAllString(text, "<strong>$1</strong>")

	// Process italic *text*
	re = regexp.MustCompile(`\*([^*]+)\*`)
	text = re.ReplaceAllString(text, "<em>$1</em>")

	return text
}

// read reads the markdown file and returns lines
func read() []string {
	file, err := os.Open("file.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return lines
}

// ParseMarkdownFile parses a markdown file and returns HTML
func ParseMarkdownFile(filename string) string {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	parser := NewMarkdownParser(lines)
	elements := parser.Parse()

	var html strings.Builder
	for _, element := range elements {
		html.WriteString(element.ToHTML())
		html.WriteString("\n")
	}

	return html.String()
}
