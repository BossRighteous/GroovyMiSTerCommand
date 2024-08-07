package display

import (
	_ "embed"
	"image"
	"image/draw"
	"strings"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

//go:embed PTSans-Regular.ttf
var ptSansRegular []byte

func parseFont(fontBytes []byte) *truetype.Font {
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		panic(err)
	}
	return font
}

var PtSansRegular *truetype.Font = parseFont(ptSansRegular)

var ScreenRect image.Rectangle = image.Rectangle{
	Min: image.Point{0, 0},
	Max: image.Point{320, 240},
}

var bgImg *image.Uniform = image.NewUniform(ColorBGR8{uint8(104), uint8(66), uint8(13)})

func ReflowText(txt string) []string {
	charsPerLine := 45
	txtLines := make([]string, 0)
	offset := 0
	end := len(txt)
	for offset < end {
		lookahead := charsPerLine
		if lookahead+offset >= end {
			lookahead = end - offset
		}
		subslice := txt[offset : offset+lookahead]
		if lookahead == charsPerLine {
			for lookahead > 0 {
				if strings.Contains(", ;.", string(subslice[lookahead-1])) {
					break
				}
				lookahead--
			}
		}
		txtLines = append(txtLines, txt[offset:offset+lookahead])
		offset += lookahead
	}
	return txtLines
}

func DrawText(text []string, rect image.Rectangle, bg *image.Uniform) *image.NRGBA {
	var (
		dpi     float64        = 72
		hinting string         = "full"
		size    float64        = 12
		spacing float64        = 1.5
		ttfont  *truetype.Font = PtSansRegular
	)

	fg := image.White
	rgba := image.NewNRGBA(rect)
	draw.Draw(rgba, rgba.Bounds(), bg, image.Point{0, 0}, draw.Src)
	c := freetype.NewContext()
	c.SetDPI(dpi)
	c.SetFont(ttfont)
	c.SetFontSize(size)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(fg)
	switch hinting {
	default:
		c.SetHinting(font.HintingNone)
	case "full":
		c.SetHinting(font.HintingFull)
	}

	// Draw the text.
	pt := freetype.Pt(10, 10+int(c.PointToFixed(size)>>6))
	for _, s := range text {
		_, _ = c.DrawString(string(s), pt)
		pt.Y += c.PointToFixed(size * spacing)
	}

	return rgba
}

func TextToBGR8(txtLines []string) *BGR8 {
	padding := []string{"", ""}
	padded := append(padding, txtLines...)
	nrgba := DrawText(padded, ScreenRect, bgImg)
	bgr := NewBGR8(ScreenRect)
	draw.Draw(bgr, ScreenRect, nrgba, image.Point{}, draw.Src)
	return bgr
}
