package sort

import (
	"testing"
	"fmt"
)

func TestQuickSort(t *testing.T) {
	arr := []int{4,6,4,5,7,3,8,0,1}
	QuickSort(arr)

	fmt.Println(arr)
}

func TestInsertSort(t *testing.T) {
	arr := []int{4,6,4,5,7,3,8,0,1}
	InsertSort(arr)

	fmt.Println(arr)
}
