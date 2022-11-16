package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"log"
	"os"

	"github.com/golang/freetype"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Pass as an argument a 2 letter string")
		os.Exit(1)
	}

	label := os.Args[1]

	// This only works for fixed width fonts btw
	fontBytes, err := ioutil.ReadFile("./noto-mono.ttf")
	if err != nil {
		log.Println(err)
		return
	}

	filename := "image.png"
	length := 100

	var palette = []color.Color{
		color.RGBA{0xff, 0xff, 0xff, 0x00}, // transparent background
		color.RGBA{0xaf, 0xff, 0xaf, 0xee}, // color 1 (purple)
		color.RGBA{0xff, 0xff, 0xff, 0xff}, // color 2 (white)
	}
	img := image.NewPaletted(image.Rect(0, 0, length, length), palette)
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
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	addCenteredLabel(img, fontBytes, float64(img.Rect.Size().X)/2, label)

	png.Encode(f, img)
	fmt.Println("Image created!")
}

func addCenteredLabel(img *image.Paletted, fontBytes []byte, fontSize float64, label string) {
	c := freetype.NewContext()
	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Println(err)
		// handle error?
		return
	}
	c.SetFont(f)
	c.SetFontSize(fontSize)
	c.SetClip(img.Bounds())
	c.SetDst(img)

	if len(label) != 2 {
		fmt.Println("The centering only works for 2 char strings")
	}
	// set letter color (this must exist in the palette btw)
	fg := image.White
	c.SetSrc(fg)
	x := float64(img.Rect.Size().X/2) - fontSize/1.7 // idk why these numbers work for any font size
	y := float64(img.Rect.Size().Y/2) + fontSize/2.6 // idk why these numbers work for any font size
	pt := freetype.Pt(int(x), int(y))
	if _, err := c.DrawString(label, pt); err != nil {
		// handle error
	}
}
