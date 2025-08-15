package request

import (
	"errors"
	"io"
	"strings"
)

type parserState int
const(
	StateInit parserState = 0
	StateDone parserState = 1
)

type Request struct{
	RequestLine RequestLine
	parser parserState
}

type RequestLine struct{
	HttpVersion string
	RequestTarget string
	Method string
}

func (r *Request) parse(p []byte) (int, error){
	switch r.parser {
	case StateDone:
		return 0, errors.New("error: trying to read data in a done state")
	case StateInit:
		rl, n, err := parseRequestLine(string(p))
		if err != nil {
			return 0, err
		}
		if n == 0 {
			return 0, nil 
		}
		r.RequestLine = *rl
		r.parser = StateDone
		return n, nil
	default:
		return 0, errors.New("error: unknown state")
	}

}
func RequestFromReader(reader io.Reader) (*Request, error){
	r := Request{parser: 0}
	readToIndex := 0
	buf := make([]byte,8)
	for r.parser != StateDone {
		if readToIndex == len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}

		n, err := reader.Read(buf[readToIndex:])
		if err != nil && err != io.EOF {
			return nil, err
		}
		readToIndex += n

		consumed, parseErr := r.parse(buf[:readToIndex])
		if parseErr != nil {
			return nil, parseErr
		}

		if consumed > 0 {
			copy(buf, buf[consumed:readToIndex])
			readToIndex -= consumed
		}
	}

	return &r, nil	
}

func parseRequestLine(data string) (*RequestLine, int ,error){
	r := strings.Split(data, "\r\n")
	if len(r) == 1{
		return nil, 0, nil
	}

	Requestline := r[0]
	n := len(Requestline) + 2
	parsed_data := strings.Split(Requestline, " ")

	if len(parsed_data) != 3{
		return nil , -1, errors.New("bad request-line request")
	}

	if parsed_data[0] != "PUT" && parsed_data[0] != "GET" && parsed_data[0] != "POST" && parsed_data[0] != "DELETE"{
		return nil , -1, errors.New("invalid method")
	} 

	parsed_data[2] = strings.Split(parsed_data[2], "/")[1]
	if parsed_data[2] != "1.1"{
		return nil , -1, errors.New("bad version of http")
	}

	rl := RequestLine{HttpVersion: parsed_data[2], RequestTarget: parsed_data[1], Method: parsed_data[0]}
	return &rl, n, nil
}