package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)
func getLinesChannel(f io.ReadCloser) <-chan string{
	c := make(chan string)
	go func(){
			line := ""
			for {
				buff := make([]byte, 8)
				n , err := f.Read(buff)
				if err != nil {
					c<- line
					close(c)
					break
				}

				data := strings.Split(string(buff), "\n")
				if len(data) == 2{
					c<- line+data[0]
					line = data[1]
					continue  
				}
				line += string(buff[:n])
			}
	}()
	return c
}

func main() {
	ls, err := net.Listen("tcp", ":8000")
	if err != nil{
		log.Fatal("Error while listening ",err)
	}
	
	for {
		conn , err := ls.Accept()
		if err != nil{
			log.Fatal("Error while making connection ",err)
		}

		c := getLinesChannel(conn)

		for line := range c{
			fmt.Println("read: ",line)
		}
		conn.Close()
	}
	
	
}