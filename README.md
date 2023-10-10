# pierre

a basic fuzzer created in Go as a learning exercise. i wrote some of my own packages for the ✨learning experience✨

## usage

The keyword `FUZZ` is used to specify which

Directory Fuzzing

```bash

go run main.go -u="http://0.0.0.0:80/FUZZ" -w /path/to/wordlist.txt

```

Fuzzing with file extensions

```bash

go run main.go -u="http://0.0.0.0:80/FUZZ" -w /path/to/wordlist.txt

```

```bash
go run main.go -u https://example.com -w wordlist.txt -e php,html,js -t 1000 -c 200 -f 200,404,403

go run main.go -u="http://0.0.0.0:80/FUZZ" -w wordlist.txt -x .php,.help,.txt,.go -fi 200,404 -fe 200,404
go run main.go -u="http://0.0.0.0:80/FUZZ" -w wordlist.txt -x .php,.help,.txt,.go -si 469 -se 469
```
