package siglas

import (
	"fmt"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"image"
	"image/color"
)

// https://stackoverflow.com/posts/54200713/revisions

func addCenteredLabel(img *image.Paletted, font *truetype.Font, fontSize float64, fgColor string, label string) {
	c := freetype.NewContext()
	c.SetFont(font)
	c.SetFontSize(fontSize)
	c.SetClip(img.Bounds())
	c.SetDst(img)

	offset := 0.0
	if len(label) != 2 {
		// log.Println("The centering only works for 2 char strings, trying to offest the width, this doesnt work in super small sizes tho")
		offset = 0.75
	}
	fgRGBAColor, err := ParseHexColor(fgColor)
	if err != nil {
		fmt.Println("unable to parse color!", fgColor)
		return
	}
	var palette = []color.Color{fgRGBAColor, fgRGBAColor, fgRGBAColor}
	fg := image.NewPaletted(image.Rect(0, 0, 2500, 2500), palette)
	c.SetSrc(fg)

	x := float64(img.Rect.Size().X/2) - fontSize/1.7 // idk why these numbers work for any font size
	y := float64(img.Rect.Size().Y/2) + fontSize/2.6 // idk why these numbers work for any font size
	pt := freetype.Pt(int(x)+int(x*offset), int(y))
	if _, err := c.DrawString(label, pt); err != nil {
		// handle error
		fmt.Println(err)
	}
}

func ParseHexColor(s string) (c color.RGBA, err error) {
	c.A = 0xff
	switch len(s) {
	case 6:
		_, err = fmt.Sscanf(s, "%02x%02x%02x", &c.R, &c.G, &c.B)
	case 4:
		_, err = fmt.Sscanf(s, "%1x%1x%1x", &c.R, &c.G, &c.B)
		// Double the hex digits:
		c.R *= 17
		c.G *= 17
		c.B *= 17
	default:
		err = fmt.Errorf("invalid length, must be 6 or 3, should not start with #")

	}
	return
}

func getImage(font *truetype.Font, bg, fg string, size int, label string) image.Image {
	length := size
	bgRGBAColor, err := ParseHexColor(bg)
	if err != nil {
		fmt.Println("Unable to parse color!", bg)
		bgRGBAColor = color.RGBA{0xaf, 0xff, 0xaf, 0xee} // color 1 (purple) default
	}
	fgRGBAColor, err := ParseHexColor(fg)
	if err != nil {
		fmt.Println("Unable to parse color!", bg)
		fgRGBAColor = color.RGBA{0xaf, 0xff, 0xff, 0xff} // White foreground by default
	}
	var palette = []color.Color{
		color.RGBA{0xff, 0xff, 0xff, 0x00}, // transparent background
		bgRGBAColor,
		fgRGBAColor,
	}
	img := image.NewPaletted(image.Rect(0, 0, length, length), palette)
	// log.Println(" Drawing background...")
	for x := 0; x < img.Rect.Size().X; x++ {
		for y := 0; y < img.Rect.Size().Y; y++ {
			// make a circle calculating if it sits within the equation of a circle,
			// maybe theres a better way of doing this
			r := img.Rect.Size().X / 2 // the radius
			if (x-r)*(x-r)+(y-r)*(y-r) <= r*r {
				img.SetColorIndex(x, y, 1)
			}
		}
	}
	// log.Println(" Drawing foreground...")
	addCenteredLabel(img, font, float64(img.Rect.Size().X)/2, fg, label)
	return img
}
