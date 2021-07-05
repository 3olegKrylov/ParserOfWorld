package main

import (
	"fmt"
	"testing"
)

type TestPair struct {
	Pair []int
}


func TestName(t *testing.T) {
	testPool := []TestPair{{[]int{1}},{[]int{1,2,3}}}
	fmt.Println(testPool)
}
