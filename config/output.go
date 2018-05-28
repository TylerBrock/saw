package config

import (
	"github.com/TylerBrock/colorjson"
	"github.com/fatih/color"
)

type OutputConfiguration struct {
	Expand         bool
	Raw            bool
	RawString      bool
	ShowStreamName bool
	Invert         bool
	NoColor        bool
}

func (c *OutputConfiguration) Formatter() *colorjson.Formatter {
	var formatter *colorjson.Formatter = colorjson.NewFormatter()

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
