// +build !solution

package main

import (
    "fmt"
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
            body, err := ioutil.ReadAll(resp.Body)
            fmt.Printf("%s\t%d\t%s\n", duration, len(body), url)
        }(url, &wg)
    }
    wg.Wait()
    totalDur := time.Since(totalStart)
    fmt.Printf("%s elapsed\n", totalDur)
}
