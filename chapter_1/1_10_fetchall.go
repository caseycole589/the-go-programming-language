package main

import(
	"fmt"
	"io"

	"net/http"
	"os"
	"time"
	"strings"

)

func main() {
	start := time.Now()
	ch := make(chan string)
	

	for _, url := range os.Args[1:] {

		file,err := os.Create(url + ".txt")
		if err != nil {
			fmt.Printf("failed creating files %s : %v\n",url,err)
			os.Exit(1)
		}
		if !strings.HasPrefix(url, "http://") {
			url = "http://" + url
		}
		go fetch(url, ch , file) //start a go subroutine
	}

	for range os.Args[1:] {
		fmt.Println(<-ch) //recieve from channel ch
	}
	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}

func fetch(url string, ch chan <- string, file *os.File) {
	
	start := time.Now()
	resp,err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprint(err)//send err to channel
		return
	}


	nbytes, err := io.Copy(file, resp.Body)
	resp.Body.Close()
	if err != nil {
		ch <- fmt.Sprintf("while reading %s: %v", url, err)//send err to main
		return
	}

	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("%.2fs %7d %s", secs, nbytes, url)
}