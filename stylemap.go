package cdtacctwebscraper

import (
	"github.com/tealeg/xlsx"
)

const (
	styleFmt_GENERAL = "general"
	styleFmt_INT     = "0"
	styleFmt_FLOAT   = "0.00"
	styleFmt_DATE    = "mm-dd-yy"
	styleFmt_ACCT    = "_(* #,##0.00_);_(* \\(#,##0.00\\);_(* \"-\"??_);_(@_)"
	styleFmt_STRING  = "@"
)

var styleMap map[string]*xlsx.Style

func init() {
	var style *xlsx.Style
	styleMap = make(map[string]*xlsx.Style)

	// General
	style = xlsx.NewStyle()
	styleMap["general"] = style

	// Table heading style
	style = xlsx.NewStyle()
	// Fill: Light blue-grey (Text 2 lighter 40%)
	style.Fill.PatternType = "solid"
	style.Fill.FgColor = "FF558ED5"
	// Font: Verdana 14, white
	style.Font.Size = 14
	style.Font.Color = "FFFFFFFF"
	// Alignment: Left
	style.Alignment.Horizontal = "left"
	styleMap["th"] = style

	// Table heading style 2
	style = xlsx.NewStyle()
	// Fill: Light orange (Accent 6 lighter 40%)
	style.Fill.PatternType = "solid"
	style.Fill.FgColor = "FFF9BE8E"
	// Font: Verdana 14, white
	style.Font.Size = 14
	style.Font.Color = "FFFFFFFF"
	// Alignment: Right
	style.Alignment.Horizontal = "left"
	styleMap["th2"] = style

	// Highlighted
	style = xlsx.NewStyle()
	style.Fill.PatternType = "solid"
	style.Fill.FgColor = "FFFFFF00"
	styleMap["hilite"] = style

	// Right justified
	style = xlsx.NewStyle()
	style.Alignment.Horizontal = "right"
	styleMap["rt"] = style

	// Thin border
	style = xlsx.NewStyle()
	style.Fill.PatternType = "solid"
	style.Fill.FgColor = "FFC6D9F1"
	style.Alignment.Horizontal = "right"
	style.Border.Left = "thin"
	style.Border.Right = "thin"
	style.Border.Bottom = "thin"
	styleMap["border"] = style
}

func getStyle(name string) (style *xlsx.Style) {
	var ok bool

	style, ok = styleMap[name]
	if !ok {
		panic("Unrecognized Style: " + name)
	}
	return
}
