// +build !solution

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
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

// CountryResponse ...
type CountryResponse struct {
	Country string `json:"country"`
	Gold    int    `json:"gold"`
	Silver  int    `json:"silver"`
	Bronze  int    `json:"bronze"`
	Total   int    `json:"total"`
}

// Person ...
type Person struct {
	Athlete      string                    `json:"athlete"`
	Country      string                    `json:"country"`
	Medals       MedalsResponse            `json:"medals"`
	MedalsByYear map[string]MedalsResponse `json:"medals_by_year"`
	sport        string
}

// Store ...
type Store struct {
	mutex         *sync.Mutex
	persons       map[string]Person
	sportPersons  map[string]map[string]Person
	sortPersons   map[string][]Person
	yearCountries map[string]map[string]CountryResponse
	sortCountries map[string][]CountryResponse
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

//PersonInfo

func (pi *Person) getGold() *int {
	return &pi.Medals.Gold
}

func (pi *Person) getSilver() *int {
	return &pi.Medals.Silver
}

func (pi *Person) getBronze() *int {
	return &pi.Medals.Bronze
}

func (pi *Person) getTotal() *int {
	return &pi.Medals.Total
}

/*  MedalsResponse */

func (pi *MedalsResponse) getGold() *int {
	return &pi.Gold
}

func (pi *MedalsResponse) getSilver() *int {
	return &pi.Silver
}

func (pi *MedalsResponse) getBronze() *int {
	return &pi.Bronze
}

func (pi *MedalsResponse) getTotal() *int {
	return &pi.Total
}

/*  CountryResponse */

func (pi *CountryResponse) getGold() *int {
	return &pi.Gold
}

func (pi *CountryResponse) getSilver() *int {
	return &pi.Silver
}

func (pi *CountryResponse) getBronze() *int {
	return &pi.Bronze
}

func (pi *CountryResponse) getTotal() *int {
	return &pi.Total
}

// IMedals ...
type IMedals interface {
	getGold() *int
	getSilver() *int
	getBronze() *int
	getTotal() *int
}

func update(val IMedals, p PersonInfo) {
	*(val.getGold()) += p.Gold
	*(val.getSilver()) += p.Silver
	*(val.getBronze()) += p.Bronze
	*(val.getTotal()) += p.Total
}

func personToMedal(p PersonInfo) MedalsResponse {
	return MedalsResponse{
		Gold:   p.Gold,
		Bronze: p.Bronze,
		Silver: p.Silver,
		Total:  p.Total,
	}
}

func getNewPerson(p PersonInfo) Person {
	return Person{
		Athlete: p.Athlete,
		Country: p.Country,
		sport:   p.Sport,
		Medals: MedalsResponse{
			Gold:   0,
			Silver: 0,
			Bronze: 0,
			Total:  0,
		},
		MedalsByYear: make(map[string]MedalsResponse),
	}
}

func getNewCountryResponse(p PersonInfo) CountryResponse {
	return CountryResponse{
		Country: p.Country,
		Gold:    0,
		Silver:  0,
		Bronze:  0,
		Total:   0,
	}
}

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
		/*
		   [ athlete info ]
		*/

		val, ok := store.persons[p.Athlete]

		if !ok {
			val = getNewPerson(p)
		}

		update(&val, p)

		val.MedalsByYear[strconv.Itoa(p.Year)] = personToMedal(p)
		store.persons[p.Athlete] = val

		/*
		   [ sport Persons info ]
		*/

		_, ok = store.sportPersons[p.Sport]
		if !ok {
			store.sportPersons[p.Sport] = make(map[string]Person)
		}
		_, ok = store.sportPersons[p.Sport][p.Athlete]
		if !ok {
			store.sportPersons[p.Sport][p.Athlete] = getNewPerson(p)
		}
		val, _ = store.sportPersons[p.Sport][p.Athlete]
		update(&val, p)
		val.MedalsByYear[strconv.Itoa(p.Year)] = personToMedal(p)

		store.sportPersons[p.Sport][p.Athlete] = val

		/*
		   [ top countries info ]
		*/

		sYear := strconv.Itoa(p.Year)
		_, ok = store.yearCountries[sYear]
		if !ok {
			store.yearCountries[sYear] = make(map[string]CountryResponse)
		}
		_, ok = store.yearCountries[sYear][p.Country]
		if !ok {
			store.yearCountries[sYear][p.Country] = getNewCountryResponse(p)
		}
		country, _ := store.yearCountries[sYear][p.Country]
		update(&country, p)
		store.yearCountries[sYear][p.Country] = country
	}

	/*
	   sort sport persons
	*/
	for sport, persons := range store.sportPersons {
		for _, p := range persons {
			store.sortPersons[sport] = append(store.sortPersons[sport], p)
		}
	}

	for key := range store.sortPersons {
		sort.Slice(store.sortPersons[key][:], func(i, j int) bool {
			switch {
			case store.sortPersons[key][i].Medals.Gold != store.sortPersons[key][j].Medals.Gold:
				return store.sortPersons[key][i].Medals.Gold > store.sortPersons[key][j].Medals.Gold
			case store.sortPersons[key][i].Medals.Silver != store.sortPersons[key][j].Medals.Silver:
				return store.sortPersons[key][i].Medals.Silver > store.sortPersons[key][j].Medals.Silver
			case store.sortPersons[key][i].Medals.Bronze != store.sortPersons[key][j].Medals.Bronze:
				return store.sortPersons[key][i].Medals.Bronze > store.sortPersons[key][j].Medals.Bronze
			default:
				return store.sortPersons[key][i].Athlete < store.sortPersons[key][j].Athlete
			}
		})
	}

	/*
	   sort countries
	*/

	for year, countries := range store.yearCountries {
		for _, c := range countries {
			store.sortCountries[year] = append(store.sortCountries[year], c)
		}
	}

	for key := range store.sortCountries {
		sort.Slice(store.sortCountries[key][:], func(i, j int) bool {
			switch {
			case store.sortCountries[key][i].Gold != store.sortCountries[key][j].Gold:
				return store.sortCountries[key][i].Gold > store.sortCountries[key][j].Gold
			case store.sortCountries[key][i].Silver != store.sortCountries[key][j].Silver:
				return store.sortCountries[key][i].Silver > store.sortCountries[key][j].Silver
			case store.sortCountries[key][i].Bronze != store.sortCountries[key][j].Bronze:
				return store.sortCountries[key][i].Bronze > store.sortCountries[key][j].Bronze
			default:
				return store.sortCountries[key][i].Country < store.sortCountries[key][j].Country
			}
		})
	}

}

// NewStore ...
func NewStore(pathToStoreFile string) *Store {
	store := &Store{
		mutex:         &sync.Mutex{},
		persons:       make(map[string]Person),
		sportPersons:  make(map[string]map[string]Person),
		sortPersons:   make(map[string][]Person),
		yearCountries: make(map[string]map[string]CountryResponse),
		sortCountries: make(map[string][]CountryResponse),
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
		if r.Method != http.MethodGet {
			s.error(w, r, http.StatusNotFound, "")
			return
		}
		r.ParseForm()
		rYear, ok := r.Form["year"]
		if !ok {
			s.error(w, r, http.StatusBadRequest, "")
			return
		}

		rLimit, ok := r.Form["limit"]
		limit := 3
		if ok {
			var err error
			limit, err = strconv.Atoi(rLimit[0])
			if err != nil || limit <= 0 {
				s.error(w, r, http.StatusBadRequest, "")
				return
			}
		}

		s.store.mutex.Lock()
		resp, ok := s.store.sortCountries[rYear[0]]
		s.store.mutex.Unlock()
		if !ok {
			err := fmt.Sprintf("year %s not found\n", rYear[0])
			s.error(w, r, http.StatusNotFound, err)
			return
		}
		if limit > len(resp) {
			limit = len(resp)
		}
		s.respond(w, r, http.StatusOK, resp[:limit])
	}
}

func (s *Server) topAthletesHandle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			s.error(w, r, http.StatusNotFound, "")
			return
		}
		r.ParseForm()
		rSport, ok := r.Form["sport"]
		if !ok {
			s.error(w, r, http.StatusBadRequest, "")
			return
		}
		sport := rSport[0]

		rLimit, ok := r.Form["limit"]
		limit := 3
		if ok {
			var err error
			limit, err = strconv.Atoi(rLimit[0])
			if err != nil || limit <= 0 {
				s.error(w, r, http.StatusBadRequest, "")
				return
			}
		}

		s.store.mutex.Lock()
		resp, ok := s.store.sortPersons[sport]
		s.store.mutex.Unlock()

		if !ok {
			err := fmt.Sprintf("sport %s not found\n", rSport)
			s.error(w, r, http.StatusNotFound, err)
			return
		}
		if limit > len(resp) {
			limit = len(resp)
		}
		s.respond(w, r, http.StatusOK, resp[:limit])
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

		s.store.mutex.Lock()
		resp, found := s.store.persons[name]
		s.store.mutex.Unlock()

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
	if len(os.Args) != 5 {
		err := fmt.Errorf("Usage: ./m -port 8080 -data /lol/kek. Need port arg you send: -->  %s", os.Args[1])
		if err != nil {
			panic(err)
		}
		return
	}

	if os.Args[1] != "-port" {
		err := fmt.Errorf("Usage: ./m -port 8080 -data /lol/kek. Need port arg you send: -->  %s", os.Args[1])
		if err != nil {
			panic(err)
		}
	}

	if os.Args[3] != "-data" {
		err := fmt.Errorf("Usage: ./m -port 8080 -data /lol/kek. Need port arg you send: -->  %s", os.Args[1])
		if err != nil {
			panic(err)
		}
	}
	srv := NewServer(os.Args[4])
	if err := http.ListenAndServe(":"+os.Args[2], srv); err != nil {
		panic(err)
	}
}
