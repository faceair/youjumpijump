package jump

import (
	"image"
	"image/color"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/nfnt/resize"
)

var basePath string

func init() {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	basePath = filepath.Dir(ex)

	if ok, _ := Exists(basePath + "/debugger"); !ok {
		os.MkdirAll(basePath+"/debugger", os.ModePerm)
	}

	os.Remove(basePath + "/debugger/debug.log")
	logFile, _ := os.OpenFile(basePath+"/debugger/debug.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	log.SetOutput(io.MultiWriter(os.Stdout, logFile))
}

func Debugger() {
	if ok, _ := Exists(basePath + "/jump.png"); ok {
		os.Rename(basePath+"/jump.png", basePath+"/debugger/"+strconv.Itoa(TimeStamp())+".png")

		files, err := ioutil.ReadDir(basePath + "/debugger/")
		if err != nil {
			panic(err)
		}

		for _, f := range files {
			fname := f.Name()
			ext := filepath.Ext(fname)
			name := fname[0 : len(fname)-len(ext)]
			if ts, err := strconv.Atoi(name); err == nil {
				if TimeStamp()-ts > 10 {
					os.Remove(basePath + "/debugger/" + fname)
				}
			}
		}
	}
}

func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func TimeStamp() int {
	return int(time.Now().UnixNano() / int64(time.Second))
}

func Distance(a, b []int) float64 {
	return math.Pow(math.Pow(float64(a[0]-b[0]), 2)+math.Pow(float64(a[1]-b[1]), 2), 0.5)
}

func getRGB(m color.Model, c color.Color) [3]int {
	if m == color.RGBAModel {
		return [3]int{int(c.(color.RGBA).R), int(c.(color.RGBA).G), int(c.(color.RGBA).B)}
	} else if m == color.RGBA64Model {
		return [3]int{int(c.(color.RGBA64).R), int(c.(color.RGBA64).G), int(c.(color.RGBA64).B)}
	} else if m == color.NRGBAModel {
		return [3]int{int(c.(color.NRGBA).R), int(c.(color.NRGBA).G), int(c.(color.NRGBA).B)}
	} else if m == color.NRGBA64Model {
		return [3]int{int(c.(color.NRGBA64).R), int(c.(color.NRGBA64).G), int(c.(color.NRGBA64).B)}
	}
	return [3]int{0, 0, 0}
}

func colorSimilar(a, b [3]int, distance float64) bool {
	return (math.Abs(float64(a[0]-b[0])) < distance) && (math.Abs(float64(a[1]-b[1])) < distance) && (math.Abs(float64(a[2]-b[2])) < distance)
}

func Find(src image.Image) ([]int, []int) {
	src = resize.Resize(720, 0, src, resize.Lanczos3)
	f, _ := os.OpenFile("jump.720.png", os.O_WRONLY|os.O_CREATE, 0600)
	png.Encode(f, src)
	f.Close()

	bounds := src.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y

	jumpCubeColor := [3]int{54, 52, 92}
	points := [][]int{}
	for y := 0; y < h; y++ {
		line := 0
		for x := 0; x < w; x++ {
			c := src.At(x, y)
			if colorSimilar(getRGB(src.ColorModel(), c), jumpCubeColor, 20) {
				line++
			} else {
				if y > 350 && x-line > 10 && line > 30 {
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
	if jumpCube[0] == 0 {
		return nil, nil
	}

	possible := [][]int{}
	for y := 0; y < h; y++ {
		line := 0
		bgColor := getRGB(src.ColorModel(), src.At(w-25, y))
		for x := 0; x < w; x++ {
			c := src.At(x, y)
			if !colorSimilar(getRGB(src.ColorModel(), c), bgColor, 5) {
				line++
			} else {
				if y > 350 && x-line > 10 && line > 35 && ((x-line/2) < (jumpCube[0]-20) || (x-line/2) > (jumpCube[0]+20)) {
					possible = append(possible, []int{x - line/2, y, line, x})
				}
				line = 0
			}
		}
	}
	if len(possible) == 0 {
		return jumpCube, nil
	}
	target := possible[0]
	for _, point := range possible {
		if point[3] > target[3] && point[1]-target[1] <= 5 {
			target = point
		}
	}
	target = []int{target[0], target[1]}

	return jumpCube, target
}
