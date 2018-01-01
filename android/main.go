package main

import (
	"fmt"
	"image/png"
	"log"
	"math"
	"os"
	"os/exec"
	"runtime/debug"
	"strconv"
	"time"

	jump "github.com/faceair/youjumpijump"
)

func main() {
	defer func() {
		jump.Debugger()
		if e := recover(); e != nil {
			log.Printf("%s: %s", e, debug.Stack())
			fmt.Print("the program has crashed, press any key to exit")
			var c string
			fmt.Scanln(&c)
		}
	}()

	var ratio float64
	var err error
	if len(os.Args) > 1 {
		ratio, err = strconv.ParseFloat(os.Args[1], 10)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Print("input jump ratio (recommend 2.04):")
		_, err = fmt.Scanln(&ratio)
		if err != nil {
			log.Fatal(err)
		}
	}
	log.Printf("now jump ratio is %f", ratio)

	for {
		jump.Debugger()

		_, err := exec.Command("/system/bin/screencap", "-p", "jump.png").Output()
		if err != nil {
			panic("screenshot failed")
		}

		inFile, err := os.Open("jump.png")
		if err != nil {
			panic(err)
		}
		src, err := png.Decode(inFile)
		if err != nil {
			panic(err)
		}
		inFile.Close()

		start, end := jump.Find(src)
		if start == nil {
			log.Print("can't find the starting point，please export the debugger directory")
			break
		} else if end == nil {
			log.Print("can't find the end point，please export the debugger directory")
			break
		}

		distance := math.Pow(math.Pow(float64(start[0]-end[0]), 2)+math.Pow(float64(start[1]-end[1]), 2), 0.5)
		log.Printf("from:%v to:%v distance:%.2f ratio:%v press:%.2fms ", start, end, distance, ratio, distance*ratio)

		scale := float64(src.Bounds().Max.X) / 720
		_, err = exec.Command("/system/bin/sh", "/system/bin/input", "swipe", strconv.FormatFloat(float64(start[0])*scale, 'f', 0, 32), strconv.FormatFloat(float64(start[1])*scale, 'f', 0, 32), strconv.FormatFloat(float64(end[0])*scale, 'f', 0, 32), strconv.FormatFloat(float64(end[1])*scale, 'f', 0, 32), strconv.Itoa(int(distance*ratio))).Output()
		if err != nil {
			panic("touch failed")
		}

		time.Sleep(time.Millisecond * 1500)
	}
}
