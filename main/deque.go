package main

import (
	"fmt"
	d "github.com/edwingeng/deque"
)

func main() {
	dq := d.NewDeque()
	dq.PushBack(100)
	dq.PushBack(200)
	dq.PushBack(300)
	for !dq.Empty() {
		fmt.Println(dq.PopFront())
	}
}
