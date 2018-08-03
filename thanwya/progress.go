package thanwya

import (
	"fmt"
	"math"
)

func printProgress(current, total int) {
	perPartition := getPerPartition(current, total)

	fmt.Print("[ ")
	for i := 1; i <= NumberOfProgressPartitions; i++ {
		if current >= min(perPartition*i, total) {
			fmt.Print("#")
		} else {
			fmt.Print("-")
		}
	}
	fmt.Printf(" ] - %0.2f%%", getPercentage(current, total))
}

func getPerPartition(current, total int) int {
	return int(math.Ceil(float64(total) / float64(NumberOfProgressPartitions)))
}

func getPercentage(a, b int) float64 {
	return float64(a) / float64(b) * 100.0
}
