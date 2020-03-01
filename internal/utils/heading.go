package utils

import "fmt"

type Heading struct {
	headingLength uint
	separator     string
}

type Align uint

const (
	Right Align = iota
	Left
)

func NewHeading(separator string, headingLength uint) *Heading {
	h := &Heading{
		headingLength: headingLength,
		separator:     separator,
	}
	return h
}

func (h *Heading) formatHeading(heading string, align Align) string {
	headingLen := len(heading)
	padding := int(h.headingLength) - headingLen - 1
	if padding < 0 {
		lastChar := headingLen + padding
		heading = fmt.Sprintf("%s ", heading[:lastChar])
	} else {
		totalWidth := padding + headingLen
		if align == Right {
			heading = fmt.Sprintf("%*s ", totalWidth, heading)
		} else {
			heading = fmt.Sprintf("%-*s ", totalWidth, heading)
		}
	}

	return heading
}

func (h *Heading) WriteAligned(heading string, text string, align Align) {
	fmt.Printf("%s%s %s\n", h.formatHeading(heading, align), h.separator, text)
}

func (h *Heading) Write(heading string, text string) {
	h.WriteAligned(heading, text, Left)
}
