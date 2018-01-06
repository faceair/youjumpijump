package jump

import (
	"image"
	"math"
)

func FindBeforeBlock(beforePic, afterPic image.Image, beforeBlock, afterBottle []int) []int {
	minX, maxX := afterBottle[0]-beforeBlock[2], afterBottle[0]+beforeBlock[2]
	minY, maxY := afterBottle[1]-(beforeBlock[1]-beforeBlock[4])*2, afterBottle[1]+(beforeBlock[1]-beforeBlock[4])*2

	block := []int{}
	for x := int(math.Max(float64(minX), 0)); x < int(math.Min(float64(maxX), float64(afterPic.Bounds().Max.X))); x++ {
		for y := int(math.Max(float64(minY), 0)); y < int(math.Min(float64(maxY), float64(afterPic.Bounds().Max.Y))); y++ {

			if IsEdge(afterPic, x, y) {
				afterLeftPixel := GetRGB(afterPic.ColorModel(), afterPic.At(x+10, y))
				beforeLeftPixel := GetRGB(beforePic.ColorModel(), beforePic.At(beforeBlock[0]-beforeBlock[2]/2+10, beforeBlock[1]))
				if ColorSimilar(afterLeftPixel, beforeLeftPixel, 20) {

					for line := beforeBlock[2] - 5; line < beforeBlock[2]+5; line++ {
						if IsEdge(afterPic, x+line, y) {
							afterRightPixel := GetRGB(afterPic.ColorModel(), afterPic.At(x+line-10, y))
							beforeRightPixel := GetRGB(beforePic.ColorModel(), beforePic.At(beforeBlock[3]-10, beforeBlock[1]))

							if ColorSimilar(afterRightPixel, beforeRightPixel, 20) && IsEdge(afterPic, x+line/2, y-(beforeBlock[1]-beforeBlock[4])) {
								block = []int{x + line/2, y}
							}
						}
					}
				}
			}
		}
	}
	return block
}
