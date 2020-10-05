// +build !solution

package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
)

func main() {
    score := make(map[string]int64)
    for _, arg := range os.Args[1:] {
        file, err := os.Open(arg)
        if err != nil {
            log.Fatal(err)
        }
        scanner := bufio.NewScanner(file)
        for scanner.Scan() {
            score[scanner.Text()]++
        }
        file.Close()
    }

    for key, val := range score {
        if val >= 2 {
            fmt.Printf("%d\t%s\n", val, key)
        }
    }
}
