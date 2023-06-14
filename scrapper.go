package main

import (
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	log "github.com/gothew/l-og"

	"golang.org/x/net/html"
)

func getLocations() {
	scrapeLocationsToJson()
}

func scrapeLocationsToJson() {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	log.Info("Getting response...")
	resp, err := client.Get("http://www.cuug.ab.ca/dewara/mandelbrot/images/images.html")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	log.Info("Parsing HTML...")

	parseHTML(resp.Body)
}

func parseHTML(body io.Reader) {
	htmlTokens := html.NewTokenizer(body)

	locFile := LocationFile{}
	currLoc := Location{}
	reachedEOF := false

	for !reachedEOF {
		tt := htmlTokens.Next()
		switch tt {
		case html.ErrorToken:
			err := htmlTokens.Err()
			if err == io.EOF {
				reachedEOF = true
			}
		case html.TextToken:
			t := htmlTokens.Token()
			parseTextToken(t.Data, &currLoc)
		case html.EndTagToken:
			t := htmlTokens.Token()
			if t.Data == "p" && currLoc.Zoom > 0 {
				locFile.Locations = append(locFile.Locations, currLoc)
			}

		}
	}

	log.Info(locFile)
}

func parseTextToken(text string, currLoc *Location) {
	splitted := strings.Split(text, " ")
	if len(splitted) != 3 {
		return
	}

	prop := strings.Replace(splitted[0], "\n", "", -1)
	switch prop {
	case "X":
		currLoc.XCenter, _ = strconv.ParseFloat(splitted[2], 64)
	case "Y":
		currLoc.YCenter, _ = strconv.ParseFloat(splitted[2], 64)
	case "R":
		rezZoom, _ := strconv.ParseFloat(splitted[2], 64)
		currLoc.Zoom = 1 / rezZoom
	}
}
