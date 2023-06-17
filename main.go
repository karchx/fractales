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
	samplesPtr := flag.Int("samples", 50, "The number of samples.")
	maxInterPtr := flag.Int("iter", 500, "The max. number of iterations.")
	hueOffsetPtr := flag.Float64("hue", 0.0, "The hsl hue offset in the range [0,1)")
	mixingPtr := flag.Bool("mixing", true, "Use linear color mixing.")
	insideBlackPtr := flag.Bool("black", true, "Paint area inside in black.")

	flag.Parse()

	imgConf = ImageConfig{
		Width:       *imgWidthPtr,
		Height:      *imgHeightPtr,
		Samples:     *samplesPtr,
		MaxInter:    *maxInterPtr,
		HueOffset:   *hueOffsetPtr,
		Mixing:      *mixingPtr,
		InsideBlack: *insideBlackPtr,
		RngGlobal:   uint64(time.Now().UnixNano()),
	}
}

func generateImagesFromLocations(locs LocationFile) {
	if _, err := os.Stat("results/" + strconv.Itoa(imgConf.MaxInter)); os.IsNotExist(err) {
		os.MkdirAll("results/"+strconv.Itoa(imgConf.MaxInter), 0755)
	}

	for index, loc := range locs.Locations {
		log.Infof("Allocating and rendering image %d\n", index+1)
		img := image.NewRGBA(image.Rect(0, 0, imgConf.Width, imgConf.Height))
		renderImage(img, loc)

		log.Infof("Encoding image %d\n", index+1)
		filename := "results/" + strconv.Itoa(imgConf.MaxInter) + "/" + strconv.Itoa(index+1)
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

	for y := 0; y < imgConf.Height; y++ {
		jobs <- y
		log.Infof("\r%d%d (%d%%)", y, imgConf.Height, int(100*(float64(y)/float64(imgConf.Height))))
	}
	log.Infof("\r%d/%[1]d (100%%)\n", imgConf.Height)
}

func renderRow(loc Location, y int, img *image.RGBA, rndLocal *uint64) {
	for x := 0; x < imgConf.Width; x++ {
		img.SetRGBA(x, y, color.RGBA{R: 0, G: 0, B: 0, A: 255})
	}
}

/*func getColorForPixel(loc Location, x int, y int, rndLocal *uint64) (uint8, uint8, uint8) {
  var r, g, b int

  for i := 0; i < imgConf.Samples
}*/
