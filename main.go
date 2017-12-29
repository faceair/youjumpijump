package main

import (
	"fmt"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
	"os/exec"
	"strconv"
	"time"
)

var jumpCubeColor = color.NRGBA{54, 52, 92, 255}

func colorSimilar(a, b color.Color, distance float64) bool {
	ra, ga, ba := a.(color.NRGBA).A, a.(color.NRGBA).G, a.(color.NRGBA).B
	rb, gb, bb := b.(color.NRGBA).A, b.(color.NRGBA).G, b.(color.NRGBA).B
	return (math.Abs(float64(ra-rb)) < distance) && (math.Abs(float64(ga-gb)) < distance) && (math.Abs(float64(ba-bb)) < distance)
}

func main() {
	var ratio float64
	fmt.Print("请输入跳跃系数:")
	_, err := fmt.Scanln(&ratio)
	if err != nil {
		log.Fatal(err)
	}

	for {
		_, err := exec.Command("adb", "shell", "screencap", "-p", "/sdcard/jump.png").Output()
		if err != nil {
			log.Fatal(err)
		}
		_, err = exec.Command("adb", "pull", "/sdcard/jump.png", ".").Output()
		if err != nil {
			log.Fatal(err)
		}

		infile, err := os.Open("jump.png")
		if err != nil {
			log.Fatal(err)
		}
		defer infile.Close()

		src, err := png.Decode(infile)
		if err != nil {
			log.Fatal(err)
		}

		bounds := src.Bounds()
		w, h := bounds.Max.X, bounds.Max.Y

		points := [][]int{}
		for y := 0; y < h; y++ {
			line := 0
			for x := 0; x < w; x++ {
				c := src.At(x, y)
				if colorSimilar(c, jumpCubeColor, 20) {
					line++
				} else {
					if y > 200 && x-line > 10 && line > 30 {
						points = append(points, []int{x - line/2, y, line})
					}
					line = 0
				}
			}
		}
		jumpCube := []int{0, 0, 0}
		for _, point := range points {
			if point[2] > jumpCube[2] {
				jumpCube = point
			}
		}
		jumpCube = []int{jumpCube[0], jumpCube[1]}

		target := []int{0, 0}
		last := []int{0, 0, 0, 0}
		for y := 0; y < h; y++ {
			line := 0
			bgColor := src.At(w-10, y)
			for x := 0; x < w; x++ {
				c := src.At(x, y)
				if !colorSimilar(c, bgColor, 10) {
					line++
				} else {
					if y > 200 && x-line > 10 && line > 35 && ((x-line/2) < (jumpCube[0]-30) || (x-line/2) > (jumpCube[0]+30)) {
						if x <= last[3] {
							target = []int{last[0], last[1]}
							break
						}
						last = []int{x - line/2, y, line, x}
					}
					line = 0
				}
			}
			if target[0] != 0 {
				break
			}
		}
		ms := int(math.Pow(math.Pow(float64(jumpCube[0]-target[0]), 2)+math.Pow(float64(jumpCube[1]-target[1]), 2), 0.5) * ratio)

		log.Printf("from:%v to:%v wait:%vms", jumpCube, target, ms)

		_, err = exec.Command("adb", "shell", "input", "swipe", "320", "410", "320", "410", strconv.Itoa(ms)).Output()
		if err != nil {
			log.Fatal(err)
		}

		time.Sleep(time.Millisecond * time.Duration(ms+420))
	}
}
