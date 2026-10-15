//go:build !solution

package main

import (
	"cmp"
	"encoding/json"
	"flag"
	"fmt"
	"maps"
	"net/http"
	"os"
	"slices"
	"strconv"
)

var winners []Winner

type Winner struct {
	Athlete string
	Country string
	Year    int
	Sport   string
	Gold    int
	Silver  int
	Bronze  int
	Total   int
}

type AthleteInfo struct {
	Athlete      string              `json:"athlete"`
	Country      string              `json:"country"`
	Medals       map[string]int      `json:"medals"`
	MedalsByYear map[string]map[string]int `json:"medals_by_year"`
}
type CountryInfo struct {
	Country string `json:"country"`
	Gold    int    `json:"gold"`
	Silver  int    `json:"silver"`
	Bronze  int    `json:"bronze"`
	Total   int    `json:"total"`
}


func athleteInfo(winners map[string][]Winner, athlete string) *AthleteInfo {
	ws, ok := winners[athlete]

	if !ok {
		return nil
	}

	medalsByYear := map[string][]int{}
	totalMedals := []int{0, 0, 0, 0}

	for _, w := range ws {
		year := fmt.Sprintf("%d", w.Year)
		
		medalsByYear[year] = []int{w.Gold, w.Silver, w.Bronze, w.Total}
		totalMedals[0] += w.Gold
		totalMedals[1] += w.Silver
		totalMedals[2] += w.Bronze
		totalMedals[3] += w.Total
	}

	info := &AthleteInfo{
		Athlete: ws[0].Athlete,
		Country: ws[0].Country,
		MedalsByYear: map[string]map[string]int{},
	}
	info.Medals = map[string]int{
		"gold":   totalMedals[0],
		"silver": totalMedals[1],
		"bronze": totalMedals[2],
		"total":  totalMedals[3],
	}

	for year, medals := range medalsByYear {
		info.MedalsByYear[year] = map[string]int{
			"gold":   medals[0],
			"silver": medals[1],
			"bronze": medals[2],
			"total":  medals[3],
		}
	}

	return info
}

func topAthletes(winners map[string][]Winner, sport string, limit int) []AthleteInfo {
	result := []AthleteInfo{}
	sportWinners := map[string] []Winner {}
	for athlete, winners := range winners {
		for _, w := range winners {
			if w.Sport == sport {
				sportWinners[athlete] = append(sportWinners[athlete], w)
			}
		}
		

	}
	for athlete := range sportWinners {
		info := athleteInfo(sportWinners, athlete)
		if info == nil {
			continue
		}
		result = append(result, *info)
	}
	slices.SortFunc(
		result,
		func(i1, i2 AthleteInfo) int {
			return -1 * cmp.Or(
				cmp.Compare(i1.Medals["gold"], i2.Medals["gold"]),
				cmp.Compare(i1.Medals["silver"], i2.Medals["silver"]),
				cmp.Compare(i1.Medals["bronze"], i2.Medals["bronze"]),
				-cmp.Compare(i1.Athlete, i2.Athlete),
			)
		},
	)
	if limit > len(result) {
		limit = len(result)
	}
	return result[:limit]
}

func topCountries(winners []Winner, year int, limit int) []CountryInfo {
	result := map[string]CountryInfo{}

	for _, winner := range winners {
		if winner.Year == year {
			info := result[winner.Country]
			info.Country = winner.Country
			info.Gold += winner.Gold
			info.Silver += winner.Silver
			info.Bronze += winner.Bronze
			info.Total += winner.Total
			result[winner.Country] = info
		}
	}
	countries := slices.Collect(maps.Values(result))

	slices.SortFunc(
		countries,
		func(c1, c2 CountryInfo) int {
			return cmp.Or(
				-cmp.Compare(c1.Gold, c2.Gold),
				-cmp.Compare(c1.Silver, c2.Silver),
				-cmp.Compare(c1.Bronze, c2.Bronze),
				cmp.Compare(c1.Country, c2.Country),
			)
		},
	)

	if limit > len(countries) {
		limit = len(countries)
	}
	return countries[:limit]
}

func athleteHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	name := query.Get("name")
	if name == "" {
		http.Error(w, "athlete not found", http.StatusNotFound)
		return
	}
	winnersMap := map[string][]Winner{}
	for _, w := range winners {
		ws := winnersMap[w.Athlete]
		winnersMap[w.Athlete] = append(ws, w)
	}

	info := athleteInfo(winnersMap, name)

	if info == nil {
		http.Error(w, "athlete not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info)

}
func topAthletesHandler(w http.ResponseWriter, r *http.Request) {
	sport := r.URL.Query().Get("sport")
	limitStr := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil && limitStr != "" {
		http.Error(w, "invalid limit", http.StatusBadRequest)
		return
	}
	if limit == 0 {
		limit = 3
	}

	winnersMap := map[string][]Winner{}
	for _, w := range winners {
		ws := winnersMap[w.Athlete]
		winnersMap[w.Athlete] = append(ws, w)
	}

	result := topAthletes(winnersMap, sport, limit)

	if len(result) == 0 {
		http.Error(w, "sport not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func topCountriesHandler(w http.ResponseWriter, r *http.Request) {
	year, err := strconv.Atoi(r.URL.Query().Get("year"))
	if err != nil {
		http.Error(w, "year not found", http.StatusNotFound)
		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil && limitStr != "" {
		http.Error(w, "invalid limit", http.StatusBadRequest)
		return
	}
	if limit == 0 {
		limit = 3
	}

	winnersMap := map[string][]Winner{}
	for _, w := range winners {
		ws := winnersMap[w.Athlete]
		winnersMap[w.Athlete] = append(ws, w)
	}

	result := topCountries(winners, year, limit)

	if len(result) == 0 {
		http.Error(w, "year not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)

}

func newMux() *http.ServeMux {
	handler := http.NewServeMux()
	handler.HandleFunc("GET /athlete-info", athleteHandler)
	handler.HandleFunc("GET /top-athletes-in-sport", topAthletesHandler)
	handler.HandleFunc("GET /top-countries-in-year", topCountriesHandler)

	return handler
}

func main() {

	port := flag.Int("port", 8080, "Port to run")
	path := flag.String("data", "", "Path to load data from")
	flag.Parse()

	file, err := os.Open(*path)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	err = json.NewDecoder(file).Decode(&winners)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	mux := newMux()
	http.ListenAndServe(fmt.Sprintf(":%d", *port), mux)
}
