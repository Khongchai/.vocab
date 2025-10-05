package super_memo

import "math"

// https://en.wikipedia.org/wiki/SuperMemo
func Sm2(grade int, repetitionNumber int, interval int, easinessFactor float64) (int, int, float64) {
	if grade >= 3 {
		switch repetitionNumber {
		case 0:
			interval = 0
		case 1:
			interval = 6
		default:
			interval = int(math.Round(float64(interval) * easinessFactor))
		}
		repetitionNumber++
	} else {
		repetitionNumber = 0
		interval = 1
	}

	newFactor := easinessFactor + (0.1-(5-float64(grade)))*(0.08+(5-float64(grade))*0.02)
	easinessFactor = math.Max(1.3, newFactor)

	return repetitionNumber, interval, easinessFactor
}
