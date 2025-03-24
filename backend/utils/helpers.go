package utils

import (
	"math/rand/v2"
)

func ShuffleSlice[T any](slice []T) {
	n := len(slice)
	for i := n - 1; i > 0; i-- {
		j := rand.IntN(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
}

func removeByValue(slice *[]string, value string) {
	for i, v := range *slice {
		if v == value {
			*slice = append((*slice)[:i], (*slice)[i+1:]...)
			return
		}
	}
}

func GetNewCol(matrix [][]string, above_used []string, total_ele []string) (newcol []string) {
	// fmt.Println("Above Used : ", above_used)
	if len(matrix) == len(matrix[0]) {
		return nil
	}
	available_ele := make([]string, len(total_ele))
	copy(available_ele, total_ele)

	for _, v := range above_used {
		removeByValue(&available_ele, v)
	}
	if len(available_ele) == 0 {
		return above_used
	}
	for _, v := range matrix[len(above_used)] {
		removeByValue(&available_ele, v)
	}

	// fmt.Println("Available eles : ", available_ele)
	if len(available_ele) == 0 {
		return nil
	}

	//insert randomizer here
	ShuffleSlice(available_ele)

	for _, v := range available_ele {
		newcol = GetNewCol(matrix, append(above_used, v), total_ele)
		if len(newcol) == 0 {
			continue
		} else {
			return newcol
		}
	}
	return nil
}
