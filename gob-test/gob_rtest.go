package main

import (
	"encoding/gob"
	"bytes"
	//"fmt"
	"fmt"
)

type Point struct {
	X int
	Y int
}

func main() {
	b := []byte{0x03,0x04,0x00,0x04}
	buf := bytes.Buffer{}
	buf.Write(b)

	var a int
	gob.NewDecoder(&buf).Decode(&a)
	fmt.Println(a)

	buf = bytes.Buffer{}
	c := []byte{0x07, 0xff , 0x82, 0x01, 0x2c, 0x01, 0x42, 0x00}
	buf.Write(c)
	var d Point
	fmt.Println(gob.NewDecoder(&buf).Decode(&d))
}
