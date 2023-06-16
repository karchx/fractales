package main

import (
	"flag"
	"image"
	"image/color"
	"image/png"
	"os"
	"runtime"
	"strconv"
	"time"

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
		Width:     *imgWidthPtr,
		Height:    *imgHeightPtr,
		RngGlobal: uint64(time.Now().UnixNano()),
	}
}

func generateImagesFromLocations(locs LocationFile) {
	if _, err := os.Stat("results/" + strconv.Itoa(500)); os.IsNotExist(err) {
		os.MkdirAll("results/"+strconv.Itoa(500), 0755)
	}

	for index, loc := range locs.Locations {
		log.Infof("Allocating and rendering image %d\n", index+1)
		img := image.NewRGBA(image.Rect(0, 0, imgConf.Width, imgConf.Height))
    renderImage(img, loc)

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

func renderImage(img *image.RGBA, loc Location) {
	jobs := make(chan int)

	for i := 0; i < runtime.NumCPU(); i++ {
		rndLocal := RandUint64(&imgConf.RngGlobal)

    go func() {
      for y := range jobs {
        renderRow(loc, y, img, &rndLocal)
      }
    }()
	}
}

func renderRow(loc Location, y int, img *image.RGBA, rndLocal *uint64) {
	for x := 0; x < imgConf.Width; x++ {
		img.SetRGBA(x, y, color.RGBA{R: 0, G: 0, B: 0, A: 255})
	}
}
