package config

import (
	"github.com/TylerBrock/colorjson"
	"github.com/fatih/color"
)

type OutputConfiguration struct {
	Raw       bool
	Pretty    bool
	Expand    bool
	Invert    bool
	RawString bool
	NoColor   bool
}

func (c *OutputConfiguration) Formatter() *colorjson.Formatter {
	formatter := colorjson.NewFormatter()

	if c.Expand {
		formatter.Indent = 4
	}

	if c.RawString {
		formatter.RawStrings = true
	}

	if c.Invert {
		formatter.KeyColor = color.New(color.FgBlack)
	}

	if c.NoColor {
		color.NoColor = true
	}

	return formatter
}
