package jump

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

var trainFile *os.File

func init() {
	trainFile, _ = os.OpenFile(basePath+"/train.ai", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
}

func NewTrain(ratio float64) *Train {
	train := &Train{
		distances:    []float64{},
		ratios:       map[float64]float64{},
		defaultRatio: ratio,
	}
	scanner := bufio.NewScanner(trainFile)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), ",")
		if len(line) == 2 {
			distance, err1 := strconv.ParseFloat(line[0], 64)
			ratio, err2 := strconv.ParseFloat(line[1], 64)
			if err1 == nil && err2 == nil {
				train.distances = append(train.distances, distance)
				train.ratios[distance] = ratio
			}
		}
	}

	return train
}

type Train struct {
	distances    []float64
	ratios       map[float64]float64
	defaultRatio float64
}

func (s *Train) Add(distance, ratio float64) {
	trainFile.Write([]byte(fmt.Sprintf("%v,%v\n", distance, ratio)))

	s.distances = append(s.distances, distance)
	s.ratios[distance] = ratio
}

func (s *Train) Find(nowDistance float64) (trainDistance, simlarRatio float64) {
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
