// +build !solution

package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "os"
)

func main() {
    for _, arg := range os.Args[1:] {
        resp, err := http.Get(arg)
        if err != nil {
            log.Fatal(err)
        }
        defer resp.Body.Close()
        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            log.Fatal(err)
        }
        fmt.Println(string(body))
    }
}
