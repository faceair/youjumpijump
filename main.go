package main

import (
	"fmt"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
	"os/exec"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/nfnt/resize"
)

var jumpCubeColor = color.NRGBA{54, 52, 92, 255}

func colorSimilar(a, b color.Color, distance float64) bool {
	ra, ga, ba := a.(color.NRGBA).A, a.(color.NRGBA).G, a.(color.NRGBA).B
	rb, gb, bb := b.(color.NRGBA).A, b.(color.NRGBA).G, b.(color.NRGBA).B
	return (math.Abs(float64(ra-rb)) < distance) && (math.Abs(float64(ga-gb)) < distance) && (math.Abs(float64(ba-bb)) < distance)
}

func main() {
	defer func() {
		if e := recover(); e != nil {
			log.Printf("%s: %s", e, debug.Stack())
			fmt.Print("程序已崩溃，请保存日志后按任意键退出\n")
			var c string
			fmt.Scanln(&c)
		}
	}()

	var ratio float64
	fmt.Print("请输入跳跃系数:")
	_, err := fmt.Scanln(&ratio)
	if err != nil {
		log.Fatal(err)
	}

	for {
		_, err := exec.Command("adb", "shell", "screencap", "-p", "/sdcard/jump.png").Output()
		if err != nil {
			panic(err)
		}
		_, err = exec.Command("adb", "pull", "/sdcard/jump.png", ".").Output()
		if err != nil {
			panic(err)
		}

		infile, err := os.Open("jump.png")
		if err != nil {
			panic(err)
		}
		defer infile.Close()

		src, err := png.Decode(infile)
		if err != nil {
			panic(err)
		}
		src = resize.Resize(720, 0, src, resize.Lanczos3)

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

		possible := [][]int{}
		for y := 0; y < h; y++ {
			line := 0
			bgColor := src.At(w-10, y)
			for x := 0; x < w; x++ {
				c := src.At(x, y)
				if !colorSimilar(c, bgColor, 10) {
					line++
				} else {
					if y > 200 && x-line > 10 && line > 35 && ((x-line/2) < (jumpCube[0]-20) || (x-line/2) > (jumpCube[0]+20)) {
						possible = append(possible, []int{x - line/2, y, line, x})
					}
					line = 0
				}
			}
		}
		target := possible[0]
		for _, point := range possible {
			if point[3] > target[3] && point[1]-target[1] <= 1 {
				target = point
			}
		}
		target = []int{target[0], target[1]}

		ms := int(math.Pow(math.Pow(float64(jumpCube[0]-target[0]), 2)+math.Pow(float64(jumpCube[1]-target[1]), 2), 0.5) * ratio)
		log.Printf("from:%v to:%v press:%vms", jumpCube, target, ms)

		_, err = exec.Command("adb", "shell", "input", "swipe", "320", "410", "320", "410", strconv.Itoa(ms)).Output()
		if err != nil {
			panic(err)
		}

		time.Sleep(time.Millisecond * 1500)
	}
}
