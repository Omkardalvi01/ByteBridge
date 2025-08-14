package main

import (
	"bytes"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", ":8000")
	if err != nil{
		log.Fatal("Error at dial", err)
	}
	defer conn.Close()

	f, err := os.Open("message.txt")
	if err != nil{
		log.Fatal("Error at Opening file", err)
	}

	data := bytes.Buffer{}
	_ , err = io.Copy(&data , f)
	if err != nil{
		log.Fatal("Error at copying file", err)
	}

	conn.Write(data.Bytes())
	
}