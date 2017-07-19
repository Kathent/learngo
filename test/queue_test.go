package test

import (
	"testing"
	"fmt"
	"queue"
)

func TestQueue(t *testing.T) {
	q := queue.New()

	q.Add(30)
	fmt.Println(q.Element())
	//q.Add(40)
	//q.Remove()
	//q.Add(100)
	//
	//for {
	//	fmt.Println(q.Element())
	//	q.Remove()
	//}
}