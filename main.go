package main

import (
	"bufio"
	"flag"
	"io"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/moul/http2curl"
	"gopkg.in/vmarkovtsev/go-lcss.v1"

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
	flag.StringVar(&scan.wordnum_include, "wi", "", "specify number of words; only allow responses with the specified number of words to be included in output")
	flag.StringVar(&scan.wordnum_exclude, "we", "", "exclude responses of a specified wordcount from output")

	//TODO: DOING THIS NOW: include first then others later
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
	post            bool
	url             string
	wordlist_file   string
	extension       string
	threads         int
	timeout         int
	filter_include  string
	filter_exclude  string
	header          string
	size_include    string
	size_exclude    string
	regex_include   string
	regex_exclude   string
	wordnum_include string
	wordnum_exclude string
}

func (s scan_config) fuzz_scan() { //post by default

	//count := 0
	output := ""

	fileread, err := os.Open(s.wordlist_file)
	handle_errors(err, "Please provide a wordlist: -w <filepath>")
	sc1 := bufio.NewScanner(fileread)
	sc1.Split(bufio.ScanLines)
	sc2 := bufio.NewScanner(fileread)
	sc2.Split(bufio.ScanLines)
	//total := count_file_lines(*sc2) * len(ext_slice)
	ext_slice := strings.Split(s.extension, ",") // TODO: put this somewhere better

	if s.post {

		//data field for post requests

	} else {

		//looops here
		for sc1.Scan() { //TODO: SHOW REDIRECTS

			for i := 0; i < len(ext_slice); i++ {

				client := &http.Client{
					Timeout: time.Duration(s.timeout) * time.Millisecond,
					CheckRedirect: func(req *http.Request, via []*http.Request) error {
						return http.ErrUseLastResponse
					},
				}

				// s.url+sc1.Text()+ext_slice[i] instead of stringsreplace
				req, err := http.NewRequest("GET", strings.Replace(s.url, "FUZZ", sc1.Text()+ext_slice[i], 1), nil) // client.Get(s.url + scanner.Text() + ext_slice[1])
				//FIXME: header must be able to be empty without index error

				if s.header != "" {
					header_slice := strings.Split(s.header, ": ") //FIXME: account for unsuccessful split
					req.Host = strings.Replace(header_slice[1], "FUZZ", sc1.Text()+ext_slice[i], 1)
					fmt.Println(header_slice[0] + "   owo    " + strings.Replace(header_slice[1], "FUZZ", sc1.Text()+ext_slice[i], 1))
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
				bodybytes, err := io.ReadAll(resp.Body)

				//if response size is not specified, include all responses!!!!
				//if size_include is specified, only include responses of that size
				//if size_exclude is specified, exclude responses of that size
				//if both are specified, only include responses of that size which are not excluded
				bodybytes_string := strconv.Itoa(count_response_bytes(bodybytes)) //keep as bytes,  but in string form
				// fmt.Println(string(bodybytes))
				fmt.Println("status code field: ", resp.StatusCode)
				//FIXME: if statement eventually gets too long, make it tidier
				//if match_int(s.size_include, bodybytes_string)
				if filter(s, bodybytes_string, resp) { //add to output if matching code

					if err != nil {
						handle_errors(err, "error counting lines")
					}

					//TODO: this formatting is painful, tidy it up
					output += sc1.Text() + ext_slice[i] + "\t\t\t[ Status: " + strconv.Itoa(resp.StatusCode) + " |" + " Size: " + strconv.Itoa(count_response_bytes(bodybytes)) + " |" + " Words: " + strconv.Itoa(count_response_words(bodybytes)) + " |" + " Lines: " + strconv.Itoa(count_response_lines(bodybytes)) + " ]\n"

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

func intersection(a, b []string) (result []string) { //get all of the common status codes, adapt later for other filtering too
	string_boo := make(map[string]bool)

	for _, thing := range a { //index = _, current element = a
		string_boo[thing] = true
	}

	for _, thing := range b {
		if _, exists := string_boo[thing]; exists {
			result = append(result, thing)
		}
	}
	return result
}

func filter(s scan_config, bodybytes_string string, resp *http.Response) bool {

	//go run main.go -u="http://0.0.0.0:80/FUZZ" -w wordlist.txt -x .sum,.go -fi 404 -fe 404 doesnt work
	status_inc_slice := strings.Split(s.filter_include, ",")
	status_ex_slice := strings.Split(s.filter_exclude, ",")

	// if s.filter_exclude and s.filter_include both include a common status code, the exclude will override the include
	lcs := lcss.LongestCommonSubstring([]byte(s.filter_include), []byte(s.filter_exclude))
	fmt.Println(string(lcs), " ", len(string(lcs)), " inc ", s.filter_include, " ex", s.filter_exclude)
	//TODO: consider splitting these strings into byte? arrays at the beginning, or with one status code per index
	if (s.filter_include != "" && s.filter_exclude != "") && len(string(lcs)) > 2 {
		// fmt.Println("filter include and filter exclude cannot include the same status code")
		// os.Exit(1)
	}

	status_match := (!strings.Contains(s.filter_exclude, strconv.Itoa(resp.StatusCode))) && strings.Contains(s.filter_include, strconv.Itoa(resp.StatusCode))

	//might need if statement for this
	bytes_match := s.size_include != "" && strings.Contains(s.size_include, bodybytes_string) || (s.size_exclude != "" && !strings.Contains(s.size_exclude, bodybytes_string)) || (s.size_include == "" && s.size_exclude == "")

	return status_match && bytes_match
}

// maybe get rid of these functions and just do it in the main function
func count_response_bytes(bodybytes []byte) int {

	return len(bodybytes)
}

func count_response_lines(bodybytes []byte) int {

	linecount := strings.Count(string(bodybytes), "\n")
	return linecount
}

func count_response_words(bodybytes []byte) int {

	word_count := strings.Count(string(bodybytes), " ")
	return word_count
}
