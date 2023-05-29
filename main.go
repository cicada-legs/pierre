package main

import (
	"bufio"
	"flag"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	// "errors"
	"fmt"
)

/*
	make some scan types, to start with

remember the "-" is implicit

first lets try to do a POST request and maybe a qwordlist
*/
func main() {

	start := time.Now()

	scan := scan_config{post: true}
	parse_flags(&scan)

	fmt.Printf(`
▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪
pierre v1.0	                                        @cicada-legs
▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪

                  ███                                       
                 ░░░                                        
       ████████  ████   ██████  ████████  ████████   ██████ 
      ░░███░░███░░███  ███░░███░░███░░███░░███░░███ ███░░███
       ░███ ░███ ░███ ░███████  ░███ ░░░  ░███ ░░░ ░███████ 
       ░███ ░███ ░███ ░███░░░   ░███      ░███     ░███░░░  
       ░███████  █████░░██████  █████     █████    ░░██████ 
       ░███░░░  ░░░░░  ░░░░░░  ░░░░░     ░░░░░      ░░░░░░  
       ░███                                                 
       █████                                                
      ░░░░░                                                 
	
▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪

❖   Url 				%s
❖   Method 				%s
❖   Wordlist 				%s
❖   Extensions 				%s
❖   Threads (Goroutines) 		%d
❖   Timeout 				%d

▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪

`, scan.url, func() string {
		if scan.post {
			return "POST"
		} else {
			return "GET"
		}
	}(), scan.wordlist_file, scan.extension, runtime.NumGoroutine(), scan.timeout)

	scan.fuzz_scan()

	duration := time.Since(start)
	fmt.Println("Total Runtime:", duration.Milliseconds(), "ms")
}

func parse_flags(scan *scan_config) {
	//
	flag.BoolVar(&scan.post, "P", false, "send HTTP request: POST by default") // TODO: make format -P GET or smth
	flag.StringVar(&scan.wordlist_file, "w", "", "provide a wordlist")         //change later??
	flag.StringVar(&scan.url, "u", "", "target url")
	flag.StringVar(&scan.extension, "x", "", "list of extensions to fuzz with: comma separated")
	flag.IntVar(&scan.threads, "t", 1, "specify the number of threads to run on")
	flag.IntVar(&scan.timeout, "c", 100, "specify the limit for when requests should timeout")

	flag.Parse()
}

type scan_config struct {
	//
	post          bool
	url           string
	wordlist_file string
	extension     string
	threads       int
	timeout       int
}

func (s scan_config) fuzz_scan() { //post by default

	ext_slice := strings.Split(s.extension, ",") // TODO: put this somewhere better

	fileread, err := os.Open(s.wordlist_file)
	handle_errors(err, "Please provide a wordlist: -w <filepath>")
	scanner := bufio.NewScanner(fileread)
	scanner.Split(bufio.ScanLines)

	if s.post {
		// resp, err := http.Post(s.url)

	} else {
		//TODO: change this to look for a word amd replace it

		for scanner.Scan() { //TODO: SHOW REDIRECTS

			for i := 0; i < len(ext_slice); i++ {

				// ctx, cancel := context.WithTimeout(context.Background(), s.timeout*time.Second)
				// defer cancel()

				// req, err := http.NewRequestWithContext(ctx, "GET", s.url+scanner.Text()+ext_slice[i], nil)

				client := &http.Client{
					Timeout: time.Duration(s.timeout) * time.Millisecond,
				}
				req, err := http.NewRequest("GET", s.url+scanner.Text()+ext_slice[1], nil) // client.Get(s.url + scanner.Text() + ext_slice[1])

				if err != nil {
					handle_errors(err, "getettet")
					return
				}
				resp, err := client.Do(req)
				if err != nil {
					fmt.Println(err)
					continue
				}

				//handle_errors(err, "second ewrrorerko")

				// client := http.Client{Timeout: s.timeout * time.Second} //TODO: change this to be the correct thing until im ready
				// fmt.Println("timeout test: ", s.timeout)
				// resp, err := client.Get(s.url + scanner.Text() + ext_slice[i])
				// if os.IsTimeout(err) {
				// 	handle_errors(err, "request timeout") //for when the request itself times out
				// 	continue                              //go to the next loop iteraion
				// } else if err != nil {
				// 	handle_errors(err, "get error") // for when the host isnt up/cant be connected to
				// // }
				defer resp.Body.Close()
				fmt.Println("/"+scanner.Text()+ext_slice[i]+"\t\t\t[ Status:", resp.StatusCode, " Size:", resp.ContentLength, "]")
				//fmt.Printf("/%s%s\t\t\t[ Status: %d  Size:%d ]\n", scanner.Text(), ext_slice[i], resp.StatusCode, resp.ContentLength)
			}
			fileread.Close()
		}
	}

}

func handle_errors(err error, msg string) {
	if err != nil {
		fmt.Println(msg)
		os.Exit(1) //change this later to be different for different errors
	}
}

func fuzz() {
	//
}
