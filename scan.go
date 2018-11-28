package main

import "fmt"
import "os"
import "net"
import "bufio"
import "strings"
import "sort"
import "log"

type Service struct {
	Name string
	Port int
	Connected bool
}

func GetService(path string) []Service {
	var services []Service

	f, err := os.Open(path); if err != nil {
		log.Fatal("os.Open: %v", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f);

	for scanner.Scan() {
		line := strings.Trim(scanner.Text(), "\r\n")
		if line != "" && line[0] == '#' { continue }
		fields := strings.Fields(line)

		service := Service {}

                if len(fields) < 2 { continue }
		service.Name = fields[0]
	        _, err := fmt.Sscanf(fields[1], "%d/tcp", &service.Port); if err == nil {
			services = append(services, service)
		}
	}

	return services
}

func Connect(ch chan Service, hostname string, service Service) {
	end := fmt.Sprintf("%s:%d", hostname, service.Port)
	conn, err := net.Dial("tcp", end); if err != nil {
		ch <- service
		return
	}

	defer conn.Close()

	service.Connected = true;

	ch <- service;
	return
}


func usage() {
	fmt.Println("./prog <host>")
	os.Exit(0)
}

type ByService []Service

func (s ByService) Len() int	{ return len(s) }
func (s ByService) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s ByService) Less(i, j int) bool { return s[i].Port < s[j].Port }

func main() {

	if len(os.Args) != 2 {
		usage()
	}

	host := os.Args[1];

	services := GetService("/etc/services");

	count := len(services)
	c := make(chan Service, count)

	fmt.Printf("Scanning %s on %d ports!\n", host, count)

	for _, s := range services {
		go Connect(c, host, s)
	}

	var results []Service;

	for i := 0; i < count; i++ {
		res := <-c
		results = append(results, res)
	}

	sort.Sort(ByService(results))

	for _, res := range results {
		if res.Connected {
			fmt.Printf("Connected %d\n", res.Port)
		}
	}

	os.Exit(0)
}
