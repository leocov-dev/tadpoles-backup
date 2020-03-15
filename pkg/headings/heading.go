package headings

import (
	"fmt"
	"github.com/fatih/color"
	"strings"
)

type Heading struct {
	headingLength uint
	separator     string
	color         *color.Color
	headingAlign  HeadingAlign
}

type HeadingAlign uint

const (
	AlignRight HeadingAlign = iota
	AlignLeft
)

type Option func(h *Heading)

func WithColor(color ...color.Attribute) Option {
	return func(h *Heading) {
		h.Color(color...)
	}
}

func AlightRight() Option {
	return func(h *Heading) {
		h.headingAlign = AlignRight
	}
}

type WriteOption uint

const (
	NoNewLine WriteOption = iota
)

func (wo WriteOption) In(opts []WriteOption) bool {
	for _, item := range opts {
		if item == wo {
			return true
		}
	}
	return false
}

func NewHeading(separator string, headingLength uint, opts ...Option) *Heading {
	h := &Heading{
		headingLength: headingLength,
		separator:     separator,
		color:         color.New(color.FgWhite),
		headingAlign:  AlignLeft,
	}

	for _, option := range opts {
		option(h)
	}
	return h
}

func (h *Heading) Copy(opts ...Option) *Heading {
	newHeading := &Heading{
		headingLength: h.headingLength,
		separator:     h.separator,
		color:         h.color,
		headingAlign:  h.headingAlign,
	}

	for _, option := range opts {
		option(newHeading)
	}
	return newHeading
}

func (h *Heading) Color(colors ...color.Attribute) {
	h.color = color.New(colors...)
}

func (h *Heading) formatHeading(heading string) string {
	headingLen := len(heading)
	padding := int(h.headingLength) - headingLen - 1
	if padding < 0 {
		lastChar := headingLen + padding
		heading = fmt.Sprintf("%s ", heading[:lastChar])
	} else {
		totalWidth := padding + headingLen
		if h.headingAlign == AlignRight {
			heading = fmt.Sprintf("%*s ", totalWidth, heading)
		} else {
			heading = fmt.Sprintf("%-*s ", totalWidth, heading)
		}
	}

	return heading
}

func (h *Heading) Write(heading string, text string, options ...WriteOption) {
	strings.TrimRight(text, "\n")
	if !NoNewLine.In(options) {
		text = text + "\n"
	}
	fmt.Printf("%s%s %s", h.color.Sprint(h.formatHeading(heading)), h.separator, text)
}
