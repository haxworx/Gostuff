package main

import "os"
import "fmt"
import "bufio"
import "io"
import "io/ioutil"
import "net/url"
import "html/template"
import "os/exec"
import "strconv"

const DISTRO_LIST_FILE = "list.txt"

type cgi struct {}

func cgi_new() cgi {
	return cgi {}
}

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

	content_length := os.Getenv("CONTENT_LENGTH")
	length, _ := strconv.Atoi(content_length)

	if length <= 0 {
		length = len(data)
	}

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

func ReadAll(path string) string {
	f, _ := os.Open(path);
	defer f.Close()
	bytes, _ := ioutil.ReadAll(f)

	return string(bytes)
}

func PhotoCode(track_url string) string {
	photocode := ""

	if len(track_url) > 0 {
		photocode = fmt.Sprintf("<iframe style=\"border: 1px solid darkred;\" width=\"220px\" height=\"220px\" src=\"%s?rel=0&autoplay=1\" frameborder=\"0\"> </iframe>", track_url)

	} else {
		photocode = "<img width=\"220px\" height=\"220px\" style=\"border: 1px solid darkred;\" src=\"/images/+icon.jpg\">"
	}

	return photocode
}

func main() {
	var data []byte
	var method string = os.Getenv("REQUEST_METHOD")

	req := cgi_new();

	switch (method) {
	case "POST":
		data = req.Post()
	case "GET":
		data = req.Get()
	}

	cmd := exec.Command("./makelist", DISTRO_LIST_FILE)

	c := make(chan int)

	go func() { cmd.Run(); c <- 1; close(c) }()

	fields := req.ParseRequest(data)
	track_url, _ := url.PathUnescape(string(fields["video"]))

	photocode := PhotoCode(track_url)
	releases := ReadAll(DISTRO_LIST_FILE)

	template_data := struct {
		PHOTOCODE template.HTML
		RELEASES  template.HTML
	} {
		PHOTOCODE: template.HTML(photocode),
		RELEASES:  template.HTML(releases),
	}

	t, _ := template.ParseFiles("templates/index.t")

	fmt.Printf("Content-type: text/html\r\n\r\n")
	t.Execute(os.Stdout, template_data)
	os.Stdout.Close()

	for i:= range c {
		fmt.Println(i)
	}
}
