package main

import (
	"flag"
	"net/http"
	"strings"

	"bufio"
	"os"

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
Pierre v1.0	                                        @cicada-legs
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
❖   Threads
❖   Timeout

▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪▪

`, scan.url, func() string {
		if scan.post {
			return "POST"
		} else {
			return "GET"
		}
	}(), scan.wordlist_file, scan.extension)

	scan.fuzz_scan()

	duration := time.Since(start)
	fmt.Println("Total Runtime:", duration.Milliseconds(), "ms")
}

func parse_flags(scan *scan_config) {
	//
	flag.BoolVar(&scan.post, "P", false, "send HTTP request: POST by default") // request type
	flag.StringVar(&scan.wordlist_file, "w", "", "provide a wordlist")         //chjange later??
	flag.StringVar(&scan.url, "u", "", "target url")
	flag.StringVar(&scan.extension, "x", "", "list of extensions to fuzz with: comma separated")

	flag.Parse()
}

type scan_config struct {
	//
	post          bool
	url           string
	wordlist_file string
	extension     string
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

		for scanner.Scan() { //SHOW REDIRECTS

			for i := 0; i < len(ext_slice); i++ {
				resp, err := http.Get(s.url + scanner.Text() + ext_slice[i])
				handle_errors(err, "get error")
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
	}
}

func fuzz() {
	//
}
