package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

func main() {
	buf := bytes.NewBuffer(nil)
	f, _ := os.Open(os.Args[1])
	io.Copy(buf, f)
	s := string(buf.Bytes())
	fmt.Println(s)
}
