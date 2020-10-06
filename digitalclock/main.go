// +build !solution

package main

import (
    "bytes"
    "fmt"
    "image"
    "image/color"
    "image/png"
    "log"
    "net/http"
    "os"
    "strconv"
    "strings"
    "time"
)

// MyTime ...
type MyTime struct {
    H int `json:"h"`
    M int `json:"m"`
    S int `json:"s"`
}

// Req ...
type Req struct {
    Time string `json:"time"`
    K    string `json:"k"`
}

// Resp ...
type Resp struct {
    Time MyTime `json:"time"`
    K    int    `json:"k"`
}

func check(r *MyTime) bool {
    if r.H < 0 || r.H > 23 {
        fmt.Println("Bad Check H")
        return false
    }

    if r.M < 0 || r.M > 59 {
        fmt.Println("Bad Check M")
        return false
    }

    if r.S < 0 || r.S > 59 {
        fmt.Println("Bad Check S")
        return false
    }
    return true
}

func getTime(t string) (MyTime, bool) {
    fmt.Printf("getTime %s \n", t)
    var res MyTime
    if t == "" {
        now := time.Now()
        res.H = now.Hour()
        res.M = now.Minute()
        res.S = now.Second()
        return res, true
    }

    tmp := strings.Split(t, ":")
    if len(tmp) != 3 {
        return res, false
    }

    var err error
    res.H, err = strconv.Atoi(tmp[0])
    if err != nil || len(tmp[0]) != 2 {
        return res, false
    }
    res.M, err = strconv.Atoi(tmp[1])
    if err != nil || len(tmp[1]) != 2 {
        return res, false
    }
    res.S, err = strconv.Atoi(tmp[2])
    if err != nil || len(tmp[2]) != 2 {
        return res, false
    }

    if !check(&res) {
        return res, false
    }
    return res, true
}

func getK(k string) (int, bool) {
    fmt.Printf("getK %s\n", k)
    if k == "" {
        return 1, true
    }

    res, err := strconv.Atoi(k)
    if err != nil || res < 1 || res > 30 {
        return res, false
    }
    return res, true
}

// SetPixel ...
func SetPixel(img *image.RGBA, pi int, pj int, k int, c color.RGBA) {
    for i := pi; i < pi+k; i++ {
        for j := pj; j < pj+k; j++ {
            img.Set(i, j, c)
        }
    }
}

// White ...
var White = color.RGBA{R: 255, G: 255, B: 255, A: 0xff}

// PrintDig ...
func PrintDig(img *image.RGBA, k int, val int, w int) {
    //for Colon
    //fmt.Printf("%v %v %v\n", val, k, w)
    dig := ``
    switch val {
    case -1:
        dig = Colon
    case 0:
        dig = Zero
    case 1:
        dig = One
    case 2:
        dig = Two
    case 3:
        dig = Three
    case 4:
        dig = Four
    case 5:
        dig = Five
    case 6:
        dig = Six
    case 7:
        dig = Seven
    case 8:
        dig = Eight
    case 9:
        dig = Nine
    default:
        panic(1)
    }

    i, j := 0, w*k
    for _, ch := range dig {
        switch ch {
        case '.':
            SetPixel(img, j, i, k, White)
            j += k
        case '1':
            SetPixel(img, j, i, k, Cyan)
            j += k
        case '\n':
            i += k
            j = w * k
        }
    }
}

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        r.ParseForm()

        structBody := &Req{
            Time: "r.FormValue()",
            K:    "r.FormValue()",
        }

        arTime, ok := r.Form["time"]
        if !ok || len(arTime) == 0 {
            structBody.Time = ""
        } else {
            if arTime[0] == "" {
                w.WriteHeader(http.StatusBadRequest)
                return
            }
            structBody.Time = arTime[0]
        }
        arK, ok := r.Form["k"]
        if !ok || len(arK) == 0 {
            structBody.K = ""
        } else {
            if arK[0] == "" {
                w.WriteHeader(http.StatusBadRequest)
                return
            }
            structBody.K = arK[0]
        }
        finalTime, ok := getTime(structBody.Time)
        if !ok {
            w.WriteHeader(http.StatusBadRequest)
            return
        }

        finalK, ok := getK(structBody.K)
        if !ok {
            w.WriteHeader(http.StatusBadRequest)
            return
        }

        w.WriteHeader(200)

        img := image.NewRGBA(image.Rect(0, 0, finalK*(8*2+4+8*2+4+8*2), finalK*12))
        PrintDig(img, finalK, finalTime.H/10, 0)
        PrintDig(img, finalK, finalTime.H%10, 0+8)
        PrintDig(img, finalK, -1, 0+8*2)
        PrintDig(img, finalK, finalTime.M/10, 0+8*2+4)
        PrintDig(img, finalK, finalTime.M%10, 0+8*2+4+8)
        PrintDig(img, finalK, -1, 0+8*2+4+8*2)
        PrintDig(img, finalK, finalTime.S/10, 0+8*2+4+8*2+4)
        PrintDig(img, finalK, finalTime.S%10, 0+8*2+4+8*2+4+8)

        f, _ := os.Create("/tmp/img.png")
        png.Encode(f, img)

        buffer := new(bytes.Buffer)
        if err := png.Encode(buffer, img); err != nil {
            log.Println("unable to encode image.")
        }

        w.Header().Set("Content-Type", "image/png")
        w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
        if _, err := w.Write(buffer.Bytes()); err != nil {
            log.Println("unable to write image.")
        }
    })

    if len(os.Args) != 3 {
        err := fmt.Errorf("Usage: ./m --port. Need two args you send: %d", len(os.Args))
        if err != nil {
            panic(err)
        }
        return
    }

    if os.Args[1] != "-port" {
        err := fmt.Errorf("Usage: ./m -port. Need port arg you send: -->  %s", os.Args[1])
        if err != nil {
            panic(err)
        }
    }
    if err := http.ListenAndServe(":"+os.Args[2], nil); err != nil {
        panic(err)
    }
}
