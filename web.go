package siglas

import (
	"image/png"
	"log"
	"net/http"
	"strconv"
)

// expected URL format:
// http://localhost:8080/?label={STRING}&bg={HEX}&fg={HEX}&size={SIZE_IN_PX}
func GetIconHandler(rw http.ResponseWriter, r *http.Request) {
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
		var err error
		sizeInt, err = strconv.Atoi(size)
		if err != nil {
			sizeInt = 255
		}
		// max to 2500
		if sizeInt > 2500 {
			sizeInt = 2500
		}
	}
	img := getImage(font, bg, fg, sizeInt, label)
	rw.Header().Set("Content-Type", "image/png")
	log.Println(" Encoding...", label, bg, fg)
	png.Encode(rw, img)
}
