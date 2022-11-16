package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
)

// https://stackoverflow.com/posts/54200713/revisions
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
	log.Println(" Drawing background...")
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
	log.Println(" Drawing foreground...")
	addCenteredLabel(img, font, float64(img.Rect.Size().X)/2, fg, label)
	return img
}
func main() {
	fontBytes, err := ioutil.ReadFile("./noto-mono.ttf")
	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Println(err)
		// handle error?
		return
	}
	if err != nil {
		log.Println(err)
		return
	}
	http.ListenAndServe(":8080", http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// get args from query params
		label := r.FormValue("label")
		bg := r.FormValue("bg")
		fg := r.FormValue("fg")
		if len(label) > 10 || len(bg) > 10 || len(fg) > 10 {
			rw.Write([]byte("Find God"))
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		size := r.FormValue("size")
		log.Println("[GET] New request!", label, bg, fg)
		sizeInt := 255
		if size != "" {
			sizeInt, err = strconv.Atoi(size)
			if err != nil {
				sizeInt = 255
			}
			// max to 2500
			if sizeInt > 2500 {
				sizeInt = 2500
			}
		}
		img := getImage(f, bg, fg, sizeInt, label)
		rw.Header().Set("Content-Type", "image/png")
		log.Println(" Encoding...", label, bg, fg)
		png.Encode(rw, img)
	}))
}

func cmd() {
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
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Println(err)
		// handle error?
		return
	}

	filename := "image.png"

	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	img := getImage(font, "FFFFFF", "000", 255, label)
	png.Encode(f, img)
	fmt.Println("Image created!")
}

func addCenteredLabel(img *image.Paletted, font *truetype.Font, fontSize float64, fgColor string, label string) {
	c := freetype.NewContext()
	c.SetFont(font)
	c.SetFontSize(fontSize)
	c.SetClip(img.Bounds())
	c.SetDst(img)

	if len(label) != 2 {
		fmt.Println("The centering only works for 2 char strings")
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
	pt := freetype.Pt(int(x), int(y))
	if _, err := c.DrawString(label, pt); err != nil {
		// handle error
		fmt.Println(err)
	}
}
