package process

import (
	"bufio"
	"encoding/csv"
	"github.com/sirupsen/logrus"
	"math"
	"os"
	"strconv"
)

func ProcessFile(filepath string) (intensity []float64, err error) {
	file, err := os.Open(filepath)
	if err != nil {
		return
	}
	reader := csv.NewReader(bufio.NewReader(file))
	records, err := reader.ReadAll()
	if err != nil {
		return
	}
	intensity = make([]float64, len(records))
	for i, record := range records {
		num, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			num = math.NaN()
		}
		intensity[i] = num
	}
	logrus.Infof("Converted %d data", len(records))
	return
}