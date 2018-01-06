package jump

import (
	"image"
	"image/color"
	"math"
)

const ExcludedHeight = 350
const ExcludedWeight = 25
const BottleWeight = 50

var BottleColor = [3]int{55, 56, 97}

func GetRGB(m color.Model, c color.Color) [3]int {
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

func ColorSimilar(a, b [3]int, distance float64) bool {
	return (math.Abs(float64(a[0]-b[0])) < distance) && (math.Abs(float64(a[1]-b[1])) < distance) && (math.Abs(float64(a[2]-b[2])) < distance)
}

func Find(pic image.Image) ([]int, []int) {
	limit := pic.Bounds().Max

	points := [][]int{}
	for gravityY := ExcludedHeight; gravityY < limit.Y; gravityY++ {
		line := 0
		for maxX := 0; maxX < limit.X-ExcludedWeight; maxX++ {
			if ColorSimilar(GetRGB(pic.ColorModel(), pic.At(maxX, gravityY)), BottleColor, 20) {
				line++
			} else {
				if line > 30 {
					points = append(points, []int{maxX - line/2, gravityY, line})
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
	for axisY := ExcludedHeight; axisY < limit.Y; axisY++ {
		line := 0
		bgColor := GetRGB(pic.ColorModel(), pic.At(limit.X-ExcludedWeight, axisY))
		for maxX := ExcludedWeight; maxX < limit.X-ExcludedWeight; maxX++ {
			if !ColorSimilar(GetRGB(pic.ColorModel(), pic.At(maxX, axisY)), bgColor, 5) {
				line++
			} else {
				axisX := maxX - line/2
				if line > 30 && (axisX < (bottle[0]-BottleWeight/2) || axisX > (bottle[0]+BottleWeight/2)) {
					points = append(points, []int{axisX, axisY, line, maxX})
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
	for topY := ExcludedHeight; topY < limit.Y; topY++ {
		bgColor := GetRGB(pic.ColorModel(), pic.At(limit.X-ExcludedWeight, topY))
		if !ColorSimilar(GetRGB(pic.ColorModel(), pic.At(block[0], topY)), bgColor, 5) {
			block = append(block, topY)
			break
		}
	}

	return bottle, block
}
