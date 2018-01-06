package jump

import (
	"image"
	"image/png"
	"math"
	"math/rand"
	"os"
	"time"
)

var randDeviation int
var randCount int

func init() {
	rand.Seed(time.Now().Unix())
}

func Distance(a, b []int) float64 {
	return math.Pow(math.Pow(float64(a[0]-b[0]), 2)+math.Pow(float64(a[1]-b[1]), 2), 0.5)
}

func GenWaitTime() int {
	waitTime := Random(1500, 5*1000)
	if Random(0, 10*1000) == 0 {
		waitTime += Random(5*1000, 10*1000)
	}
	if Random(0, 60) == 0 {
		waitTime += Random(50*1000, 60*1000)
	}
	return waitTime
}

func GenRandDeviation(block []int) ([]int, int) {
	randDeviation = 0
	if block[2] > 142 {
		if randCount == 0 {
			randDeviation = Random(18, 28)
			if Random(0, 2) == 0 {
				randDeviation *= -1
			}
		}
		randCount++
		if randCount == 2 {
			randCount = 0
		}
	}
	block[0] += randDeviation
	return block, randDeviation
}

func Random(min, max int) int {
	return rand.Intn(max-min) + min
}

func SavePNG(fileName string, pic image.Image) {
	f, _ := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0600)
	png.Encode(f, pic)
	f.Close()
}

func OpenPNG(fileName string) (image.Image, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return png.Decode(f)
}
