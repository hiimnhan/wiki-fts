package common

import "fmt"

func Pop(arr *[]int) int {
	x := (*arr)[len(*arr)-1]
	*arr = (*arr)[:len(*arr)-1]

	return x
}

func Shift(slice *[]int) (int, error) {
	if len(*slice) == 0 {
		return 0, fmt.Errorf("cannot shift from an empty slice")
	}
	firstElement := (*slice)[0]
	*slice = (*slice)[1:]
	return firstElement, nil
}

func RemoveElement(arr []int, ele int) []int {
	var index int
	for i, e := range arr {
		if e == ele {
			index = i
		}
	}

	return append(arr[:index], arr[index+1:]...)
}

// func Intersection(a, b []int) []int {
// 	maxLen := len(a)
// 	if len(b) > maxLen {
// 		maxLen = len(b)
// 	}
// 	r := make([]int, 0, maxLen)
// 	var i, j int
// 	for i < len(a) && j < len(b) {
// 		if a[i] < b[j] {
// 			i++
// 		} else if a[i] > b[j] {
// 			j++
// 		} else {
// 			r = append(r, a[i])
// 			i++
// 			j++
// 		}
// 	}
// 	return r
// }

func Intersect(arrays [][]int) []int {
	// Create a map to count the frequency of each number
	freq := make(map[int]int)
	for _, slice := range arrays {
		for _, num := range slice {
			freq[num]++
		}
	}

	// Find the numbers that appear in all slices
	var result []int
	for num, count := range freq {
		if count == len(arrays) {
			result = append(result, num)
		}
	}

	return result
}
