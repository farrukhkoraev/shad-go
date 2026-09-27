//go:build !solution

package main

import (
	"fmt"
	"time"

	"golang.org/x/tour/tree"
)

func Index[T comparable](s []T, x T) int {
	for i, v := range s {
		// v and x are type T, which has the comparable
		// constraint, so we can use == here.
		if v == x {
			return i
		}
	}
	return -1
}

type List[T any] struct {
	next *List[T]
	val  T
}

func FromArray[T any](ts []T) *List[T] {
	var head *List[T] = nil
	var cur *List[T] = head
	
	for _, t := range(ts) {
		if cur == nil {
			cur = &List[T]{nil, t}
			head = cur
		} else {
			cur.next = &List[T]{nil, t}
			cur = cur.next
		}
				
	}
	return head
	
}

func Length[T any] (hd *List[T]) int {
	var cur = hd
	var sum = 0
	
	for cur != nil {
		sum += 1
		cur = cur.next
	}
	return sum
}
func Tail[T any] (l List[T]) List[T] {
	return *(l.next)
}

func say (s string) {
	for _ = range(5) {
		time.Sleep(100 * time.Millisecond )
		fmt.Println(s)
	}
}

func sum (s []int, c chan int) {
	sum := 0
	for _, v := range s {
		sum += v
	}
	c <- sum 

}

func fibonacci(n int, c chan int) {
	x, y := 0, 1

	for _ = range n {
		c <- x
		x, y = y, x + y
		
	}
	close(c)
}

func fibonacci2(c, quit chan int) {
	x, y := 0, 1

	for {
		select {
		case c <- x:
			x, y = y, y + x 
		case <- quit:
			fmt.Println("quit")
		    return 
		}
				
	}
}

func runIndex () {
	// Index works on a slice of ints 
	si := []int{10, 20, 15, -10}
	fmt.Println(Index(si, 15))
	
	
	// Index also works on a slice of strings
	ss := []string{"foo", "bar", "baz"}
	fmt.Println(Index(ss, "hello"))
}

func runListTest() {
	ss := []string{"foo", "bar", "baz"}
	l := FromArray(ss)
	
	fmt.Println(Length(l))	
	for p := l; p != nil; {
		fmt.Printf("{val: %v}", p.val)
		p = p.next
	}	
}

func main() {


	// go say("hello")
	// say("world")

	// s := []int {1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	// 
	// c := make(chan int)
	// go sum(s[:len(s)/2], c)
	// go sum(s[len(s)/2:], c)
	// 
	// x, y := <- c, <- c
	// 
	// fmt.Println(x, y, x + y )

	// c := make(chan int, 10)
	// go fibonacci(cap(c), c)
	// for i := range c {
	// 	fmt.Println(i)
	// }

	// c := make(chan int)
	// quit := make(chan int)
	// go func() {
	// 	for i := 0; i < 10; i++ {
	// 		fmt.Println(<-c)
	// 	}
	// 	quit <2- 0
	// }()
	// fibonacci2(c, quit)

	// ch := make(chan int, 20)
	// go Walk(tree.New(1), ch)
	// 
    //  for _ = range 10 {
	// 	 fmt.Println(<-ch)
	// 	 
	// }

	fmt.Print(Same(tree.New(1), tree.New(2)))
}

func Walk(tree *tree.Tree, ch chan int) {

	if tree == nil {
		return
	}
	Walk(tree.Left, ch)
	ch <- tree.Value
	Walk(tree.Right, ch)
}

func Same (t1 *tree.Tree, t2 *tree.Tree) bool{
	var ch1, ch2 = make(chan int), make(chan int)
	go Walk(t1, ch1)
	go Walk(t2, ch2)

	for _ = range 10 {
		if <-ch1 != <-ch2 { return false}
	}
	return true
	
}
