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
