// +build !solution

package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
    "strconv"
    "sync"
)

/*
   [STRUCTS]
   ...
*/

// PersonInfo ...
type PersonInfo struct {
    Athlete string `json:"athlete"`
    Age     int    `json:"age"`
    Country string `json:"country"`
    Year    int    `json:"year"`
    Date    string `json:"date"`
    Sport   string `json:"sport"`
    Gold    int    `json:"gold"`
    Silver  int    `json:"silver"`
    Bronze  int    `json:"bronze"`
    Total   int    `json:"total"`
}

// MedalsResponse ...
type MedalsResponse struct {
    Gold   int `json:"gold"`
    Silver int `json:"silver"`
    Bronze int `json:"bronze"`
    Total  int `json:"total"`
}

// Person ...
type Person struct {
    Athlete      string                    `json:"athlete"`
    Country      string                    `json:"country"`
    Medals       MedalsResponse            `json:"medals"`
    MedalsByYear map[string]MedalsResponse `json:"medals_by_year"`
}

// Store ...
type Store struct {
    mutex   *sync.Mutex
    persons map[string]Person
}

// Server ...
type Server struct {
    mux   *http.ServeMux
    store *Store
}

/*

   [STORE]
    ...
*/

// Init ...
func (store *Store) Init(pathToStoreFile string) {
    dat, err := ioutil.ReadFile(pathToStoreFile)
    if err != nil {
        panic(err)
    }
    persons := make([]PersonInfo, 0)
    if err := json.Unmarshal(dat, &persons); err != nil {
        panic(err)
    }

    for _, p := range persons {
        val, ok := store.persons[p.Athlete]

        if !ok {
            val = Person{
                Athlete:      p.Athlete,
                Country:      p.Country,
                MedalsByYear: make(map[string]MedalsResponse),
            }
        }
        val.Medals.Gold += p.Gold
        val.Medals.Bronze += p.Bronze
        val.Medals.Silver += p.Silver
        val.Medals.Total += p.Total

        val.MedalsByYear[strconv.Itoa(p.Year)] = MedalsResponse{
            Gold:   p.Gold,
            Bronze: p.Bronze,
            Silver: p.Silver,
            Total:  p.Total,
        }
        store.persons[p.Athlete] = val
    }
}

// NewStore ...
func NewStore(pathToStoreFile string) *Store {
    store := &Store{
        mutex:   &sync.Mutex{},
        persons: make(map[string]Person),
    }
    store.Init(pathToStoreFile)
    return store
}

/*
   [Server]
   ServeHTTP
   NewServer
   Handles
   error
   respond
*/

func (s *Server) error(w http.ResponseWriter, r *http.Request, status int, err string) {
    w.Header().Add("Content-Type", "text/plain")
    w.WriteHeader(status)
    w.Write([]byte(err))
}

func (s *Server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
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

func (s *Server) topCountryHandle() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        s.respond(w, r, http.StatusOK, nil)
    }
}

func (s *Server) topAthletesHandle() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        s.respond(w, r, http.StatusOK, nil)
    }
}

func (s *Server) athleteInfoHandle() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodGet {
            s.error(w, r, http.StatusNotFound, "")
            return
        }
        r.ParseForm()
        name := r.Form.Get("name")
        found := false
        resp, found := s.store.persons[name]
        if !found {
            err := fmt.Sprintf("athlete %s not found\n", name)
            s.error(w, r, http.StatusNotFound, err)
            return
        }
        s.respond(w, r, http.StatusOK, resp)
    }
}

// Init ...
func (s *Server) Init() {
    s.mux.HandleFunc("/athlete-info", s.athleteInfoHandle())
    s.mux.HandleFunc("/top-athletes-in-sport", s.topAthletesHandle())
    s.mux.HandleFunc("/top-countries-in-year", s.topCountryHandle())
}

// NewServer ...
func NewServer(pathToStoreFile string) *Server {
    srv := &Server{
        mux:   http.NewServeMux(),
        store: NewStore(pathToStoreFile),
    }
    srv.Init()
    return srv
}

// ServeHTTP
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    s.mux.ServeHTTP(w, r)
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

    const pathToStoreFile = "testdata/olympicWinners.json"
    srv := NewServer(pathToStoreFile)
    if err := http.ListenAndServe(":"+os.Args[2], srv); err != nil {
        panic(err)
    }

}
