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
	//
	flag.BoolVar(&scan.post, "P", false, "send HTTP request: POST by default") // TODO: make format -P GET or smth
	flag.StringVar(&scan.wordlist_file, "w", "", "provide a wordlist")         //change later??
	flag.StringVar(&scan.url, "u", "", "target url")
	flag.StringVar(&scan.extension, "x", "", "list of extensions to fuzz with: comma separated")
	flag.IntVar(&scan.threads, "th", 1, "specify the number of threads to run on")
	flag.IntVar(&scan.timeout, "to", 100, "specify the limit for when requests should timeout")
	flag.StringVar(&scan.filter_codes_string, "fs", "200,204,301,302,307,400,401,403,404,405,500", "specify which status codes to be included in output") //TODO: add more codes!!!

	flag.Parse()
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
	total := count_file_lines(*sc2) * len(ext_slice)

	if s.post {
		// resp, err := http.Post(s.url)

	} else {
		//TODO: change this to look for a word amd replace it

		for sc1.Scan() { //TODO: SHOW REDIRECTS

			for i := 0; i < len(ext_slice); i++ {
				count++
				loading_animation(count, total, "<", ">", 100, "=") //TODO: this isnt printing for some reason

				client := &http.Client{
					Timeout: time.Duration(s.timeout) * time.Millisecond,
				}
				req, err := http.NewRequest("GET", s.url+sc1.Text()+ext_slice[i], nil) // client.Get(s.url + scanner.Text() + ext_slice[1])

				if err != nil {
					handle_errors(err, "getettet")
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
					output += output + "/" + sc1.Text() + ext_slice[i] + "\t\t\t[ Status: " + strconv.Itoa(resp.StatusCode) + " Size: " + strconv.FormatInt(resp.ContentLength, 10) + " ]\n"
					// fmt.Printf("/%s%s\t\t\t[ Status: %d  Size:%d ]\n", scanner.Text(), ext_slice[i], resp.StatusCode, resp.ContentLength)
					//fmt.Println(output)

				}

			}
			fileread.Close()
		}

	}
	fmt.Println(output)
}

func handle_errors(err error, msg string) {
	if err != nil {
		fmt.Printf(msg)
		os.Exit(1) //change this later to be different for different errors
	}
}

func count_file_lines(sc bufio.Scanner) int {
	lines := 0
	for sc.Scan() {
		lines++
	}
	return lines
}

func loading_animation(iterate int, total int, prefix string, suffix string, length int, fill string) {
	percentage := float64(iterate) / float64(total)
	bar_progress := int(length * iterate / total)
	end := "="

	if iterate == total {
		end = "]"
	}

	bar := strings.Repeat(fill, bar_progress) + end + strings.Repeat("-", bar_progress)
	fmt.Printf("\r%s [%s] %f%% %s", prefix, bar, percentage, suffix)
	if iterate == total {
		fmt.Println()
	}
}
