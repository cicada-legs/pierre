package main

import (
	"bufio"
	"flag"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/moul/http2curl"

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
	flag.IntVar(&scan.timeout, "to", 500, "specify the limit for when requests should timeout")
	flag.StringVar(&scan.filter_include, "fi", "200,204,301,302,307,400,401,403,405,500", "filter include; specify which status codes to be included in output")
	flag.StringVar(&scan.filter_exclude, "fe", "404", "filter exclude; specify which status codes to be excluded in output. Overrides filter include")
	flag.StringVar(&scan.size_include, "si", "", "specify response size; only allow responses of the specified size to be included in output")
	flag.StringVar(&scan.size_exclude, "se", "", "exclude responses of a specified size from output")

	//TODO: later
	flag.StringVar(&scan.regex_include, "ri", "", "specify a regex pattern to be included in output")
	flag.StringVar(&scan.regex_exclude, "re", "", "specify a regex pattern to be excluded from output")

	//TODO: add more codes!!!
	flag.StringVar(&scan.header, "H", "", "specify a header to be sent with the request. Example: Host: FUZZ.example.com")

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
	post           bool
	url            string
	wordlist_file  string
	extension      string
	threads        int
	timeout        int
	filter_include string
	filter_exclude string
	header         string
	size_include   string
	size_exclude   string
	regex_include  string
	regex_exclude  string
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
				if strings.Contains(s.filter_include, strconv.Itoa(resp.StatusCode)) { //add to output if matching code
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
				//TODO: now!!!! the header isnt being put into the request properly
				//ALSO ITERATING THROUGH HEADERS!!!

				if s.header != "" {
					header_slice := strings.Split(s.header, ": ") //TODO: account for unsuccessful split
					req.Host = strings.Replace(header_slice[1], "FUZZ", sc1.Text()+ext_slice[i], 1)
					fmt.Println(header_slice[0] + "   owo    " + strings.Replace(header_slice[1], "FUZZ", sc1.Text()+ext_slice[i], 1))
					// req.Header.Add("Host", "FUZZ.0.0.0.0")
				}
				command, _ := http2curl.GetCurlCommand(req)
				fmt.Println(command)

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
				body_bytes, err := ioutil.ReadAll(resp.Body)

				//if response size is not specified, include all responses!!!!
				//TODO: CURRENTLY DOING THIS
				//if size_include is specified, only include responses that match the size
				//same for not matching
				//if size_include and size_sxclude both contain a number like 400, then size_sxclude should take precedence and 400 should be excludeed (skip the if statement)

				if (s.size_include == "" || ((!strings.Contains(s.size_exclude, strconv.Itoa(count_response_bytes(body_bytes)))) && strings.Contains(s.size_include, strconv.Itoa(count_response_bytes(body_bytes))))) &&
					(s.size_exclude == "" || !strings.Contains(s.size_exclude, strconv.Itoa(count_response_bytes(body_bytes)))) ||
					((!strings.Contains(s.filter_exclude, strconv.Itoa(resp.StatusCode))) && strings.Contains(s.filter_include, strconv.Itoa(resp.StatusCode))) { //add to output if matching code

					if err != nil {
						handle_errors(err, "error counting lines")
					}

					//FIXME: this formatting is painful, tidy it up
					output += sc1.Text() + ext_slice[i] + "\t\t\t[ Status: " + strconv.Itoa(resp.StatusCode) + " |" + " Size: " + strconv.Itoa(count_response_bytes(body_bytes)) + " |" + " Words: " + strconv.Itoa(count_response_words(body_bytes)) + " |" + " Lines: " + strconv.Itoa(count_response_lines(body_bytes)) + " ]\n"

					//DELETE BELOW LATER
					// read and print entire response
					// body, err := httputil.DumpResponse(resp, true)
					// if err != nil {
					// 	handle_errors(err, "error reading body")
					// }
					// output += "~~~~~~~~~~~~~~~~~~~~~~~~~~~\n" + string(body) + "~~~~~~~~~~~~~~~~~~~~~~~~~~~\n" + "\n"

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

// maybe get rid of these functions and just do it in the main function
func count_response_bytes(body_bytes []byte) int {

	return len(body_bytes)
}

func count_response_lines(body_bytes []byte) int {

	linecount := strings.Count(string(body_bytes), "\n")
	return linecount
}

func count_response_words(body_bytes []byte) int {

	word_count := strings.Count(string(body_bytes), " ")
	return word_count
}
