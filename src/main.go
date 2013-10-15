package main

import (
	"flag"
	"fmt"
	"os"
	// ""
	"./rayTracer"
	"image/png"
	"log"
	"runtime/pprof"
)

var fileName *string = flag.String("f", "", "-f \"fileName\"")
var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {

	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	if *fileName == "" {
		fmt.Println("Usage ", os.Args[0], "-f \"fileName of config file\"")
		return
	}
	fmt.Println("Reading config")
	fmt.Println(*fileName)

	s, err := rayTracer.NewScene(*fileName)
	if err != nil {
		println("error", err.Error())
		return
	}
	img := rayTracer.Draw(s)
	f, err := os.Create("../out/" + s.Name + ".png")
	defer f.Close()
	if err != nil {
		println("error", err.Error())
		return
	}

	png.Encode(f, img)
	// for x := 0 ; x < s.CM.L.Bounds().Dx(); x++ {
	// 	for y := 0 ; y < s.CM.L.Bounds().Dy(); y++ {
	// 		i := s.CM.L
	// 		println(i.Rect.Min.Y,i.Rect.Min.X)
	// 		r := i.Pix[(y-i.Rect.Min.Y)*i.Stride + (x-i.Rect.Min.X)*8]
	// 		g := i.Pix[(y-i.Rect.Min.Y)*i.Stride + (x-i.Rect.Min.X)*8 +1]
	// 		b := i.Pix[(y-i.Rect.Min.Y)*i.Stride + (x-i.Rect.Min.X)*8 +2]

	// 	println(r,g,b)
	// }
	// }

	// println(r,g,b,s.CM.L.Bounds().Dx(),s.CM.L.Bounds().Dy())

}
