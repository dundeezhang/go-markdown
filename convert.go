package main

import "strconv"

// ToHTML methods for all markdown elements

func (h Heading) ToHTML() string {
	return "<h" + strconv.Itoa(h.Level) + ">" + h.Text + "</h" + strconv.Itoa(h.Level) + ">"
}

func (p Paragraph) ToHTML() string {
	return "<p>" + p.Text + "</p>"
}

func (l List) ToHTML() string {
	tag := "ul"
	if l.Ordered {
		tag = "ol"
	}
	html := "<" + tag + ">"
	for _, item := range l.Items {
		html += "<li>" + item + "</li>"
	}
	html += "</" + tag + ">"
	return html
}

func (l Link) ToHTML() string {
	return "<a href=\"" + l.URL + "\">" + l.Text + "</a>"
}

func (c CodeBlock) ToHTML() string {
	if c.Language != "" {
		return "<pre><code class=\"language-" + c.Language + "\">" + c.Code + "</code></pre>"
	}
	return "<pre><code>" + c.Code + "</code></pre>"
}

func (i InlineCode) ToHTML() string {
	return "<code>" + i.Code + "</code>"
}
