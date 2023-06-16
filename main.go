package main

import (
	"flag"
	"image"
	"image/png"
	"os"
	"strconv"

	log "github.com/gothew/l-og"
)

var imgConf ImageConfig

func main() {
	parseImageConfigArgs()
	generateImagesFromLocations(getLocations())
}

func parseImageConfigArgs() {
	imgWidthPtr := flag.Int("width", 1920, "The width of the image in pixels.")
	imgHeightPtr := flag.Int("height", 1024, "The height of the image in pixels.")

	flag.Parse()

	imgConf = ImageConfig{
		Width:  *imgWidthPtr,
		Height: *imgHeightPtr,
	}
}

func generateImagesFromLocations(locs LocationFile) {
	if _, err := os.Stat("results/" + strconv.Itoa(500)); os.IsNotExist(err) {
		os.Mkdir("results/"+strconv.Itoa(500), 0755)
	}

	for index, loc := range locs.Locations {
		log.Infof("Allocating and rendering image %d\n", index+1)
		img := image.NewRGBA(image.Rect(0, 0, imgConf.Width, imgConf.Height))

		log.Info(loc)
		log.Infof("Encoding image %d\n", index+1)
		filename := "results/" + strconv.Itoa(500) + "/" + strconv.Itoa(index+1)
		f, err := os.Create(filename + ".png")
		if err != nil {
			panic(err)
		}

		err = png.Encode(f, img)
		if err != nil {
			panic(err)
		}
	}
}
