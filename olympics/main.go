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

// Person ...
type Person struct {
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

// Store ...
type Store struct {
    mutex  *sync.Mutex
    person []Person
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
    if err := json.Unmarshal(dat, &store.person); err != nil {
        panic(err)
    }
}

// NewStore ...
func NewStore(pathToStoreFile string) *Store {
    store := &Store{
        mutex:  &sync.Mutex{},
        person: make([]Person, 0),
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

    type MedalsResponse struct {
        Gold   int `json:"gold"`
        Silver int `json:"silver"`
        Bronze int `json:"bronze"`
        Total  int `json:"total"`
    }
    type Response struct {
        Athlete      string                    `json:"athlete"`
        Country      string                    `json:"country"`
        Medals       MedalsResponse            `json:"medals"`
        MedalsByYear map[string]MedalsResponse `json:"medals_by_year"`
    }

    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodGet {
            s.error(w, r, http.StatusNotFound, "")
            return
        }
        r.ParseForm()
        name := r.Form.Get("name")
        found := false
        resp := &Response{
            Athlete:      "",
            Country:      "",
            Medals:       MedalsResponse{},
            MedalsByYear: make(map[string]MedalsResponse),
        }
        for _, person := range s.store.person {
            if person.Athlete != name {
                continue
            }
            if !found {
                resp.Athlete = person.Athlete
                resp.Country = person.Country
            }
            resp.Medals.Gold += person.Gold
            resp.Medals.Bronze += person.Bronze
            resp.Medals.Silver += person.Silver
            resp.Medals.Total += person.Total

            resp.MedalsByYear[strconv.Itoa(person.Year)] = MedalsResponse{
                Gold:   person.Gold,
                Bronze: person.Bronze,
                Silver: person.Silver,
                Total:  person.Total,
            }
            found = true
        }
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
