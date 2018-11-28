package main

import "os"
import "fmt"
import "bufio"
import "io"

type cgi struct {}

func (h cgi) Post() []byte {
	in := bufio.NewReader(os.Stdin)
	var data[]byte
	for {
		b, err := in.ReadByte(); if err == io.EOF {
			break;
		}

		if err != nil {
			os.Exit(1 << 1)
		}

		data = append(data, b)
	}

	return data
}

func (h cgi) Get() []byte {
	return []byte(os.Getenv("QUERY_STRING"))
}

func (h cgi) ParseRequest(data []byte) map[string][]byte {
	fields  := make(map[string][]byte)
	var index int = 0
	length := len(data)

	key_start, key_end := 0, 0
	val_start, val_end := 0, 0

	for index < length {
		if data[index] == '=' {
			key_end = index
			val_start = index + 1
		}

		if data[index] == '&' || index == length-1 {
			key := string(data[key_start:key_end])
			key_start = index + 1
			if index == length - 1 {
				val_end = index +  1
			} else {
				val_end = index
			}
			fields[key] = data[val_start:val_end]
		}

		index++
	}
	return fields
}

func main() {
	var incoming []byte
	var method string = os.Getenv("REQUEST_METHOD")

	req := cgi {};

	switch (method) {
		case "POST":
			incoming = req.Post()
		case "GET":
			incoming = req.Get()
		default:
			os.Exit(0 << 1)
	}

	fmt.Printf("Content-type: text/plain\r\n\r\n")
	fields := req.ParseRequest(incoming)

	for key, value := range fields {
		fmt.Printf("key: %s value: %s\n", key, string(value))
	}

	os.Exit(0)
}
