package internal

import "image/color"

type Color struct {
	hexCode string
	rgba    color.RGBA
}

func (c Color) HexCode() string {
	return c.hexCode
}

func (c Color) RGBA() color.RGBA {
	return c.rgba
}

var (
	Yellow = Color{
		hexCode: "FFFF00",
		rgba: color.RGBA{
			R: 255,
			G: 255,
			B: 0,
		},
	}
	Green = Color{
		hexCode: "00FF00",
		rgba: color.RGBA{
			R: 0,
			G: 255,
			B: 0,
		},
	}
	Red = Color{
		hexCode: "FF0000",
		rgba: color.RGBA{
			R: 255,
			G: 0,
			B: 0,
		},
	}
	Brown = Color{
		hexCode: "663300",
		rgba: color.RGBA{
			R: 102,
			G: 51,
			B: 0,
		},
	}
)
