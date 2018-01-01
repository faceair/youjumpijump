package jump

import (
	"encoding/gob"
	"math"
	"os"
)

func NewSimilar(ratio float64) *Similar {
	similar := new(Similar)
	f, err := os.Open("similar.ai")
	if err == nil {
		decoder := gob.NewDecoder(f)
		err = decoder.Decode(similar)
	}
	f.Close()
	similar.defaultRatio = ratio
	return similar
}

type Similar struct {
	distances    []float64
	ratios       map[float64]float64
	defaultRatio float64
}

func (s *Similar) Add(distance, ratio float64) {
	s.distances = append(s.distances, distance)
	s.ratios[distance] = ratio
}

func (s *Similar) Find(nowDistance float64) (similarDistance, simlarRatio float64) {
	sumR := 0.0
	sumD := 0.0
	count := 0.0

	for _, distance := range s.distances {
		if math.Abs(nowDistance-distance) < 10 {
			count++
			sumD += distance
			sumR += s.ratios[distance]
		}
	}
	if count < 3 {
		return 0, s.defaultRatio
	}

	return sumD / count, sumR / count
}

func (s *Similar) Save() {
	f, err := os.Create("similar.ai")
	if err == nil {
		encoder := gob.NewEncoder(f)
		encoder.Encode(s)
	}
	f.Close()
}
