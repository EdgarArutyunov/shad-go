// +build !solution

package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	totalStart := time.Now()

	for _, url := range os.Args[1:] {
		wg.Add(1)
		go func(url string, wg *sync.WaitGroup) {
			defer wg.Done()
			start := time.Now()
			resp, err := http.Get(url)
			duration := time.Since(start)
			if err != nil {
				fmt.Println(err)
				return
			}
			nbytes, err := io.Copy(ioutil.Discard, resp.Body)
			defer resp.Body.Close()
			fmt.Printf("%s\t%d\t%s\n", duration, nbytes, url)
		}(url, &wg)
	}
	wg.Wait()
	totalDur := time.Since(totalStart)
	fmt.Printf("%s elapsed\n", totalDur)
}
