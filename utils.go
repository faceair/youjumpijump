package jump

import (
	"image"
	"image/png"
	"math"
	"math/rand"
	"os"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func Distance(a, b []int) float64 {
	return math.Pow(math.Pow(float64(a[0]-b[0]), 2)+math.Pow(float64(a[1]-b[1]), 2), 0.5)
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
