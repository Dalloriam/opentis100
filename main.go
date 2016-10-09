package main

import (
	"fmt"

	"github.com/dalloriam/opentis/tis100"
)

func main() {
	comp := tis100.New("myTIS")
	fmt.Println(comp.Name)
}
