package main

import (
	"bufio"
	"fmt"
	"strings"
)

func main() {
	// column()
	convert()

}

var source = `1行目
2行目
3行目
`

func column() {
	reader := bufio.NewReader(strings.NewReader(source))
	for {
		line, err := reader.ReadString('\n')
		fmt.Printf("%#v\n", line)
		if err != nil {
			break
		}
	}
}

var source2 = "123 1.234 1.0e4 test"

func convert() {
	reader := strings.NewReader(source2)
	var i int
	var f, g float64
	var s string
	fmt.Fscan(reader, &i, &f, &g, &s)
	fmt.Printf("i=%#v f=%#v g=%#v s=%#v\n", i, f, g, s)
}
