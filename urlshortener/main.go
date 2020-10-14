// +build !solution

package main

import (
    "crypto/md5"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
    "sync"
)

type server struct {
    mux   *http.ServeMux
    m     *sync.Mutex
    store map[string]string
}

func (s *server) error(w http.ResponseWriter, r *http.Request, status int, err string) {
    w.Header().Add("Content-Type", "text/plain")
    w.WriteHeader(status)
    w.Write([]byte(err))
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
    w.Header().Add("Content-Type", "application/json")
    w.WriteHeader(code)

    if data != nil {
        js, err := json.Marshal(data)
        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            return
        }
        w.Write(js)
    }
}

func (s *server) shortenHandle() http.HandlerFunc {
    type Request struct {
        URL *string `json:"url"`
    }

    type Response struct {
        URL string `json:"url"`
        Key string `json:"key"`
    }

    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            s.error(w, r, http.StatusBadRequest, "invalid request")
            return
        }
        req := &Request{}

        data, err := ioutil.ReadAll(r.Body)
        if err != nil {
            s.error(w, r, http.StatusBadRequest, "invalid request")
            return
        }
        if err := json.Unmarshal(data, req); err != nil {
            s.error(w, r, http.StatusBadRequest, "invalid request")
            return
        }

        if req.URL == nil {
            s.error(w, r, http.StatusBadRequest, "invalid request")
            return
        }

        byteKey := md5.Sum([]byte(*req.URL))
        resp := &Response{
            URL: *req.URL,
            Key: fmt.Sprintf("%x", byteKey),
        }

        s.m.Lock()
        s.store[resp.Key] = resp.URL
        s.m.Unlock()

        s.respond(w, r, http.StatusOK, resp)
    }
}

func (s *server) goKeyHandle() http.HandlerFunc {
    const prefixURL = "/go/"

    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodGet {
            s.error(w, r, http.StatusBadRequest, "invalid request")
            return
        }
        suffix := r.URL.RequestURI()
        if len(suffix) < len(prefixURL) {
            s.error(w, r, http.StatusInternalServerError, "")
            return
        }

        suffix = suffix[len(prefixURL):]

        s.m.Lock()
        resURL, ok := s.store[suffix]
        s.m.Unlock()

        if !ok {
            s.error(w, r, http.StatusNotFound, "key not found")
            return
        }
        w.Header().Add("Location", resURL)

        w.WriteHeader(http.StatusFound)
    }
}

func (s *server) configureMux() {
    s.mux.HandleFunc("/shorten", s.shortenHandle())
    s.mux.HandleFunc("/go/", s.goKeyHandle())
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    s.mux.ServeHTTP(w, r)
}

func newServer() *server {
    s := &server{
        mux:   http.NewServeMux(),
        store: make(map[string]string),
        m:     &sync.Mutex{},
    }
    s.configureMux()
    return s
}

func main() {
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
    srv := newServer()
    if err := http.ListenAndServe(":"+os.Args[2], srv); err != nil {
        panic(err)
    }

}
