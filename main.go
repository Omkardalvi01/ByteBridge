package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"github.com/Omkardalvi01/ByteBridge/internal/request"
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

		r , err := request.RequestFromReader(conn)
		if err != nil{
			log.Fatal("Error while reading ",err)
		}
		
		fmt.Println("Method:",r.RequestLine.Method)
		fmt.Println("Version:",r.RequestLine.HttpVersion)
		fmt.Println("Addr:",r.RequestLine.RequestTarget)
		conn.Close()
	}
	
	
}