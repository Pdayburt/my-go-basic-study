package syntax

import (
	"fmt"
	"testing"
)

func TestSelectChannel(t *testing.T) {

	ch1 := make(chan int, 1)
	ch2 := make(chan int, 2)
	ch2 <- 100
	ch1 <- 200

	select {
	case val := <-ch1:
		t.Log("val1:", val)
		val2 := <-ch2
		t.Log("val2:", val2)
	case val := <-ch2:
		t.Log("val2:", val)
		val1 := <-ch1
		t.Log("val1:", val1)
	}

}

func TestChannel(t *testing.T) {
	//声明一个chan，但没有初始化。读写都会崩溃
	//var ch chan int64

	//初始化一个容量为0的chan
	//ch1 := make(chan int64)

	//初始化一个容量为2的chan
	ch2 := make(chan int64, 3)
	//关闭ch2 不能写但是能读，读出来都是类型的零值
	defer close(ch2)
	ch2 <- 10086
	ch2 <- 10010
	i, ok := <-ch2
	if !ok {
		//ch2 已经被人关闭了
	}
	fmt.Println(ok, i)
}
