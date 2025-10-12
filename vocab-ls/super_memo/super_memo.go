package super_memo

import "math"

const (
	MemoBlackout            int = 0 // 0 - complete blackout
	MemoIncorrectRemembered int = 1 // 1 - incorrect response; the correct one remembered
	MemoIncorrectEasy       int = 2 // 2 - incorrect response; the correct one seemed easy to recall
	MemoCorrectHard         int = 3 // 3 - correct response recalled with serious difficulty
	MemoCorrectHesitation   int = 4 // 4 - correct response after a hesitation
	MemoPerfect             int = 5 // 5 - perfect response
)

const InitialEasinessFactor = 2.5

// https://en.wikipedia.org/wiki/SuperMemo
func Sm2(grade int, repetitionNumber int, interval float64, easinessFactor float64) (int, float64, float64) {
	if grade >= 3 {
		switch repetitionNumber {
		case 0:
			interval = 1
		case 1:
			interval = 6
		default:
			interval = math.Round(interval * easinessFactor)
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
