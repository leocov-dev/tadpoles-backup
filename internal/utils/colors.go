package utils

import (
	"github.com/fatih/color"
	"tadpoles-backup/pkg/headings"
)

var (
	headingBase = headings.NewHeading(":", 18)
	WriteMain   = headingBase.Copy(
		headings.WithColor(color.FgHiYellow),
	).Write
	WriteSub = headingBase.Copy(
		headings.WithAlignRight(),
		headings.WithColor(color.FgHiBlue),
	).Write
	WriteInfo = headingBase.Copy(
		headings.WithColor(color.FgHiMagenta),
	).Write
	WriteError = headingBase.Copy(
		headings.WithColor(color.FgHiRed),
	).Write
	WriteErrorSub = headingBase.Copy(
		headings.WithAlignRight(),
		headings.WithColor(color.FgHiRed),
	)
)
