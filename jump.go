package jump

import (
	"image"
	"image/color"
	"math"
	"os"

	"github.com/nfnt/resize"
)

var ExcludedHeight = 350
var ExcludedWeight = 25
var BottleColor = [3]int{55, 56, 97}

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

func Find(pic image.Image) ([]int, []int) {
	pic = resize.Resize(720, 0, pic, resize.Lanczos3)

	if len(os.Getenv("DEBUG")) > 0 {
		go SavePNG("jump.720.png", pic)
	}

	bounds := pic.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y

	points := [][]int{}
	for y := 0; y < h; y++ {
		line := 0
		for x := 0; x < w; x++ {
			c := pic.At(x, y)
			if colorSimilar(getRGB(pic.ColorModel(), c), BottleColor, 20) {
				line++
			} else {
				if y > ExcludedHeight && x-line > ExcludedWeight && line > 30 {
					points = append(points, []int{x - line/2, y, line})
				}
				line = 0
			}
		}
	}
	bottle := []int{0, 0, 0}
	for _, point := range points {
		if point[2] >= bottle[2] && point[1] > bottle[1] {
			bottle = point
		}
	}
	bottle = []int{bottle[0], bottle[1]}
	if bottle[0] == 0 {
		return nil, nil
	}

	points = [][]int{}
	for y := 0; y < h; y++ {
		line := 0
		bgColor := getRGB(pic.ColorModel(), pic.At(w-ExcludedWeight, y))
		for x := 0; x < w; x++ {
			c := pic.At(x, y)
			if !colorSimilar(getRGB(pic.ColorModel(), c), bgColor, 5) {
				line++
			} else {
				if y > ExcludedHeight && x-line > ExcludedWeight && line > 30 &&
					((x-line/2) < (bottle[0]-20) || (x-line/2) > (bottle[0]+20)) {
					points = append(points, []int{x - line/2, y, line, x})
				}
				line = 0
			}
		}
	}
	if len(points) == 0 {
		return bottle, nil
	}
	block := points[0]
	for _, point := range points {
		if point[3] > block[3] && point[1]-block[1] < 30 && math.Abs(float64(block[0]-point[0])) < 2 {
			block = point
		}
	}
	block = []int{block[0], block[1]}

	return bottle, block
}
