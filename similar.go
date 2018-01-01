package jump

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

var similarFile *os.File

func init() {
	similarFile, _ = os.OpenFile(basePath+"/similar.ai", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
}

func NewSimilar(ratio float64) *Similar {
	similar := &Similar{
		distances:    []float64{},
		ratios:       map[float64]float64{},
		defaultRatio: ratio,
	}
	scanner := bufio.NewScanner(similarFile)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), ",")
		if len(line) == 2 {
			distance, err1 := strconv.ParseFloat(line[0], 64)
			ratio, err2 := strconv.ParseFloat(line[1], 64)
			if err1 == nil && err2 == nil {
				similar.Add(distance, ratio)
			}
		}
	}

	return similar
}

type Similar struct {
	distances    []float64
	ratios       map[float64]float64
	defaultRatio float64
}

func (s *Similar) Add(distance, ratio float64) {
	similarFile.Write([]byte(fmt.Sprintf("%v,%v\n", distance, ratio)))

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
