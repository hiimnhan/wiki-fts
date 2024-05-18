package common

func Pop(arr *[]int) int {
	x := (*arr)[len(*arr)-1]
	*arr = (*arr)[:len(*arr)-1]

	return x
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
