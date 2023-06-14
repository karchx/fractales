package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	log "github.com/gothew/l-og"

	"golang.org/x/net/html"
)

func getLocations() LocationFile {
	file, err := os.ReadFile(LocationFileJson)
	if err != nil {
		if os.IsNotExist(err) {
			scrapeLocationsToJson()
		} else {
			panic(err)
		}
	}

	locs := LocationFile{}
	_ = json.Unmarshal([]byte(file), &locs)

	zoom1Fractal := Location{
		XCenter: -0.75,
		YCenter: 0,
		Zoom:    1,
	}

	locs.Locations = append(locs.Locations, zoom1Fractal)

	log.Infof("Found %v locations.", len(locs.Locations))
	return locs
}

func scrapeLocationsToJson() {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	log.Info("Getting response...")
	resp, err := client.Get(UrlScapper)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	log.Info("Parsing HTML...")

	locFile := parseHTML(resp.Body)

	log.Info("Writing location data to JSON...")
	res, err := json.MarshalIndent(locFile, "", " ")
	if err != nil {
		log.Fatal(err)
	} else {
		_ = os.WriteFile(LocationFileJson, res, 0644)
	}

	log.Info("Scraping location data successfull.")
}

func parseHTML(body io.Reader) LocationFile {
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
	return locFile
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
