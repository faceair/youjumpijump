package jump

import (
	"image"
	"math"
)

const QE = 100.0

func IsEdge(pic image.Image, x, y int) bool {
	var sum float64

	pixelO := pic.At(x, y)
	x2, y2, z2, _ := pixelO.RGBA()
	for _, offsetX := range []int{-1, 0, 1} {
		for _, offsetY := range []int{-1, 0, 1} {
			if offsetX == offsetY {
				break
			}
			pixelN := pic.At(x+offsetX, y+offsetY)
			x1, y1, z1, _ := pixelN.RGBA()

			xSqr := (x1 - x2) * (x1 - x2)
			ySqr := (y1 - y2) * (y1 - y2)
			zSqr := (z1 - z2) * (z1 - z2)
			mySqr := float64(xSqr + ySqr + zSqr)
			dist := math.Sqrt(mySqr)
			sum += dist
		}
	}
	return sum/8 > 65536/QE
}
