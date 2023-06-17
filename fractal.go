package main

type Location struct {
	XCenter float64
	YCenter float64
	Zoom    float64
}

type LocationFile struct {
	Locations []Location
}

type ImageConfig struct {
	Width       int
	Height      int
	Samples     int
	MaxInter    int
	HueOffset   float64
	Mixing      bool
	InsideBlack bool
	RngGlobal   uint64
}
