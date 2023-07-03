package main

import (
	"bufio"
	"flag"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	// "errors"
	"fmt"
)

/*
	example command:
	go run main.go -u https://www.google.com -w wordlist.txt -x php,html,txt -th 10 -to 100 -fs 200,204,301,302,307,400,401,403,404,405,500

	/TODO: add a keyword which is replaced by the wordlist


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

❖   URL 				%s
❖   Method 				%s
❖   Wordlist 				%s
❖   Extensions 				%s
❖   Threads (Goroutines) 		%d
❖   Timeout (seconds) 			%d

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

	flag.BoolVar(&scan.post, "P", false, "send HTTP request: GET by default") // TODO: make format -P GET or smth
	flag.StringVar(&scan.wordlist_file, "w", "", "provide a wordlist")        //change later??
	flag.StringVar(&scan.url, "u", "", "target url")
	flag.StringVar(&scan.extension, "x", "", "list of extensions to fuzz with: comma separated")
	flag.IntVar(&scan.threads, "th", 1, "specify the number of threads to run on")
	flag.IntVar(&scan.timeout, "to", 100, "specify the limit for when requests should timeout")
	flag.StringVar(&scan.filter_codes_string, "fi", "200,204,301,302,307,400,401,403,404,405,500", "filter include; specify which status codes to be included in output")
	//TODO: add more codes!!!
	flag.StringVar(&scan.header, "H", "", "specify a header to be sent with the request. Example: Host: FUZZ.example.com")
	//flag.StringVar(&scan.filter_codes_string, "filter-status", "200,204,301,302,307,400,401,403,404,405,500", "specify which status codes to be included in output")

	flag.Parse()

	//TODO: this condition is bad, change it later

	if scan.url == "" || !((!strings.Contains(scan.url, "FUZZ") && strings.Contains(scan.header, "FUZZ")) ||
		(strings.Contains(scan.url, "FUZZ") && !strings.Contains(scan.header, "FUZZ"))) {
		fmt.Println("Please provide \"FUZZ\" in the URL where you want to fuzz\n Example: http://example.com/FUZZ")
		flag.PrintDefaults()
		os.Exit(1)
	}
}

type scan_config struct {
	//
	post                bool
	url                 string
	wordlist_file       string
	extension           string
	threads             int
	timeout             int
	filter_codes_string string
	header              string
}

func (s scan_config) fuzz_scan() { //post by default

	count := 0
	output := ""
	ext_slice := strings.Split(s.extension, ",") // TODO: put this somewhere better

	fileread, err := os.Open(s.wordlist_file)
	handle_errors(err, "Please provide a wordlist: -w <filepath>")
	sc1 := bufio.NewScanner(fileread)
	sc1.Split(bufio.ScanLines)
	sc2 := bufio.NewScanner(fileread)
	sc2.Split(bufio.ScanLines)
	//total := count_file_lines(*sc2) * len(ext_slice)

	if s.post {

		for sc1.Scan() {

			for i := 0; i < len(ext_slice); i++ {
				count++
				// TODO: check below for errors from copilot

				client := &http.Client{
					Timeout: time.Duration(s.timeout) * time.Millisecond,
				}
				req, err := http.NewRequest("POST", s.url+sc1.Text()+ext_slice[i], nil) // client.Get(s.url + scanner.Text() + ext_slice[1])

				if err != nil {
					handle_errors(err, "ow")
					return
				}
				resp, err := client.Do(req)
				if err != nil {
					// fmt.Println(err) /this is for testing only
					//connection timeout, what to output at the end of the scan?
					continue
				}

				defer resp.Body.Close()
				// fmt.Println("code", resp.StatusCode)
				if strings.Contains(s.filter_codes_string, strconv.Itoa(resp.StatusCode)) { //add to output if matching code
					output += sc1.Text() + ext_slice[i] + "\n"
				}
			}

		}

	} else { //TODO: adapt output for subdomain fuzzing

		for sc1.Scan() { //TODO: SHOW REDIRECTS

			for i := 0; i < len(ext_slice); i++ {

				client := &http.Client{
					Timeout: time.Duration(s.timeout) * time.Millisecond,
					CheckRedirect: func(req *http.Request, via []*http.Request) error {
						return http.ErrUseLastResponse
						//TODO: RIGHT NOWWWWWW< MAKE SURE TO PRINT THIS WAWAWAW
					},
				}

				// s.url+sc1.Text()+ext_slice[i] instead of stringsreplace
				req, err := http.NewRequest("GET", strings.Replace(s.url, "FUZZ", sc1.Text()+ext_slice[i], 1), nil) // client.Get(s.url + scanner.Text() + ext_slice[1])
				//FIXME: header must be able to be empty without index error

				if s.header != "" {
					header_slice := strings.Split(s.header, ": ") //TODO: account for unsuccessful split
					req.Header.Add(header_slice[0], header_slice[1])
					// req.Header.Add("Host", "FUZZ.0.0.0.0")
				}

				if err != nil {
					handle_errors(err, "GET error: "+err.Error())
					return
				}

				resp, err := client.Do(req)

				if err != nil {
					// fmt.Println(err) /this is for testing only
					//connection timeout, what to output at the end of the scan?
					continue
				}

				defer resp.Body.Close()

				if strings.Contains(s.filter_codes_string, strconv.Itoa(resp.StatusCode)) { //add to output if matching code
					output += sc1.Text() + ext_slice[i] + "\t\t\t[ Status: " + strconv.Itoa(resp.StatusCode) + " Size: " + strconv.FormatInt(resp.ContentLength, 10) + " ]\n"
					//DELETE BELOW LATER
					//read and print entire response
					// body, err := httputil.DumpResponse(resp, true)
					// if err != nil {
					// 	handle_errors(err, "error reading body")
					// }
					// output += "~~~~~~~~~~~~~~~~~~~~~~~~~~~\n" + string(body) + "~~~~~~~~~~~~~~~~~~~~~~~~~~~\n" + "\n"

					//another

					//convert resp.Location() to string
					// redirect_url, err := resp.Location()
					// if err != nil {
					// 	handle_errors(err, err.Error())
					// }
					// output += redirect_url.String() + "\n"

					//another

					// for k, v := range resp.Header {
					// 	fmt.Print(k)
					// 	fmt.Print(" : ")
					// 	fmt.Println(v)
					// }
				}

			}
		}

	}
	fmt.Println(output)
}

func handle_errors(err error, msg string) {
	if err != nil {
		fmt.Println(msg)
		os.Exit(1) //change this later to be different for different errors
	}
}
