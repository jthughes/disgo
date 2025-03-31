package main

import "fmt"

func clamp(i, min, max int) (int, error) {
	if min > max {
		return 0, fmt.Errorf("min larger than max")
	}
	if i < min {
		return min, nil
	}
	if i > max {
		return max, nil
	}
	return i, nil
}
