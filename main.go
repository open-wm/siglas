package siglas

import (
	"embed"
	"fmt"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
)

//go:embed noto-mono.ttf
var f embed.FS

var font *truetype.Font

func main() {
	SampleServer()
}

func ReadDefaultFont() error {
	fontBytes, err := f.ReadFile("noto-mono.ttf")
	if err != nil {
		log.Println("Unable to read font", err)
		// handle error?
		return err
	}
	font, err = freetype.ParseFont(fontBytes)
	if err != nil {
		log.Println("Unable to parse font", err, len(fontBytes))
		return err
	}
	return nil
}

func SampleServer() {
	if err := ReadDefaultFont(); err != nil {
		log.Println("Error, couldnt start the server", err)
		return
	}
	http.HandleFunc("/", GetIconHandler)
	log.Println("Starting server in port :8080")
	http.ListenAndServe(":8080", nil)
}

func cli() {
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
