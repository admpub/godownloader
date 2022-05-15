package monitor_test

import (
	"fmt"
	"testing"
)

func TestDefer(t *testing.T) {
	defer func() {
		fmt.Println("============1")
	}()
	defer func() {
		fmt.Println("============2")
	}()
	//panic("~")
	//output:
	//============2
	//============1
}
