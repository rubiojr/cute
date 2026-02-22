package cairoops

import (
	"strconv"
	"strings"

	"github.com/jwijenbergh/purego"
)

var (
	cairoSetSourceRGBA  func(uintptr, float64, float64, float64, float64)
	cairoRectangle      func(uintptr, float64, float64, float64, float64)
	cairoFill           func(uintptr)
	cairoMoveTo         func(uintptr, float64, float64)
	cairoLineTo         func(uintptr, float64, float64)
	cairoSetLineWidth   func(uintptr, float64)
	cairoStroke         func(uintptr)
	cairoSelectFontFace func(uintptr, string, int, int)
	cairoSetFontSize    func(uintptr, float64)
	cairoShowText       func(uintptr, string)
)

func init() {
	lib, err := purego.Dlopen("libcairo.so.2", purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		panic(err)
	}

	purego.RegisterLibFunc(&cairoSetSourceRGBA, lib, "cairo_set_source_rgba")
	purego.RegisterLibFunc(&cairoRectangle, lib, "cairo_rectangle")
	purego.RegisterLibFunc(&cairoFill, lib, "cairo_fill")
	purego.RegisterLibFunc(&cairoMoveTo, lib, "cairo_move_to")
	purego.RegisterLibFunc(&cairoLineTo, lib, "cairo_line_to")
	purego.RegisterLibFunc(&cairoSetLineWidth, lib, "cairo_set_line_width")
	purego.RegisterLibFunc(&cairoStroke, lib, "cairo_stroke")
	purego.RegisterLibFunc(&cairoSelectFontFace, lib, "cairo_select_font_face")
	purego.RegisterLibFunc(&cairoSetFontSize, lib, "cairo_set_font_size")
	purego.RegisterLibFunc(&cairoShowText, lib, "cairo_show_text")
}

func parseHexByte(v string) (float64, bool) {
	n, err := strconv.ParseUint(v, 16, 8)
	if err != nil {
		return 0, false
	}
	return float64(n) / 255.0, true
}

func colorRGBA(color string) (float64, float64, float64, float64) {
	s := strings.TrimSpace(strings.ToLower(color))
	if strings.HasPrefix(s, "#") {
		hex := s[1:]
		switch len(hex) {
		case 3:
			r, rok := parseHexByte(strings.Repeat(string(hex[0]), 2))
			g, gok := parseHexByte(strings.Repeat(string(hex[1]), 2))
			b, bok := parseHexByte(strings.Repeat(string(hex[2]), 2))
			if rok && gok && bok {
				return r, g, b, 1
			}
		case 6:
			r, rok := parseHexByte(hex[0:2])
			g, gok := parseHexByte(hex[2:4])
			b, bok := parseHexByte(hex[4:6])
			if rok && gok && bok {
				return r, g, b, 1
			}
		case 8:
			r, rok := parseHexByte(hex[0:2])
			g, gok := parseHexByte(hex[2:4])
			b, bok := parseHexByte(hex[4:6])
			a, aok := parseHexByte(hex[6:8])
			if rok && gok && bok && aok {
				return r, g, b, a
			}
		}
	}

	switch s {
	case "white":
		return 1, 1, 1, 1
	case "red":
		return 1, 0, 0, 1
	case "green":
		return 0, 1, 0, 1
	case "blue":
		return 0, 0, 1, 1
	}
	return 0, 0, 0, 1
}

func Fill(ctx uintptr, width float64, height float64, color string) {
	if ctx == 0 || width <= 0 || height <= 0 {
		return
	}
	r, g, b, a := colorRGBA(color)
	cairoSetSourceRGBA(ctx, r, g, b, a)
	cairoRectangle(ctx, 0, 0, width, height)
	cairoFill(ctx)
}

func Line(ctx uintptr, x1 float64, y1 float64, x2 float64, y2 float64, color string) {
	if ctx == 0 {
		return
	}
	r, g, b, a := colorRGBA(color)
	cairoSetSourceRGBA(ctx, r, g, b, a)
	cairoSetLineWidth(ctx, 1.0)
	cairoMoveTo(ctx, x1, y1)
	cairoLineTo(ctx, x2, y2)
	cairoStroke(ctx)
}

func Text(ctx uintptr, x float64, y float64, value string, color string) {
	if ctx == 0 {
		return
	}
	r, g, b, a := colorRGBA(color)
	cairoSetSourceRGBA(ctx, r, g, b, a)
	cairoSelectFontFace(ctx, "Sans", 0, 0)
	cairoSetFontSize(ctx, 13.0)
	cairoMoveTo(ctx, x, y)
	cairoShowText(ctx, value)
}

func Rect(ctx uintptr, x float64, y float64, w float64, h float64, color string) {
	if ctx == 0 || w <= 0 || h <= 0 {
		return
	}
	r, g, b, a := colorRGBA(color)
	cairoSetSourceRGBA(ctx, r, g, b, a)
	cairoRectangle(ctx, x, y, w, h)
	cairoFill(ctx)
}
