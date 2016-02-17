package main

import (
	"log"
	"net/http"

	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	"regexp"
	"strconv"
	"strings"
)

type Params struct {
	Port string `envconfig:"PORT" required:"true"`
	DSN  string `required:"true"`
}

var countryMap = map[string]bool{}

func main() {
	// parse params
	var params Params
	if err := envconfig.Process("fsmigrations", &params); err != nil {
		log.Fatal("Error parsing parameters", err)
	}
	log.Printf("Params %#v\n", params)

	// init countryMap
	for _, c := range countries {
		countryMap[strings.ToLower(c)] = true
	}

	// init db
	db, err := sqlx.Open("mysql", params.DSN)
	if err != nil {
		log.Fatal("Error opening db", err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Println("Error pinging db", err)
	}

	// init server
	log.Println("Start server")
	e := echo.New()

	e.Use(mw.Logger())
	e.Use(mw.Recover())

	e.Get("/", func(c *echo.Context) error {
		return c.String(http.StatusOK, "Ok")
	})

	e.Get("/immigrations", func(c *echo.Context) error {
		place := c.Query("place")
		year := c.Query("year")
		filter := c.Query("filter")

		placeCounts, err := getImmigrations(db, place, year)
		if err != nil {
			return err
		}

		result := asResult(placeCounts, filter)

		b, err := json.Marshal(result)
		if err != nil {
			return err
		}
		return c.String(http.StatusOK, removeCountQuotes(string(b)))
	})

	e.Get("/emigrations", func(c *echo.Context) error {
		place := c.Query("place")
		year := c.Query("year")
		filter := c.Query("filter")

		placeCounts, err := getEmigrations(db, place, year)
		if err != nil {
			return err
		}

		result := asResult(placeCounts, filter)

		b, err := json.Marshal(result)
		if err != nil {
			return err
		}
		return c.String(http.StatusOK, removeCountQuotes(string(b)))
	})

	e.Run(":" + params.Port)
}

type DbSelecter interface {
	Select(dest interface{}, query string, args ...interface{}) error
}

type PlaceCount struct {
	Place string
	Count int
}

type ImmigrationPlaceCount struct {
	Place string `db:"place_from"`
	Count int    `db:"count"`
}

type EmigrationPlaceCount struct {
	Place string `db:"place_to"`
	Count int    `db:"count"`
}

func getImmigrations(db DbSelecter, place, year string) ([]PlaceCount, error) {
	placeCounts := []ImmigrationPlaceCount{}
	err := db.Select(&placeCounts, "SELECT place_from, count FROM immigrations WHERE place_to = ? and year = ?", place, year)
	if err != nil {
		return nil, err
	}

	result := make([]PlaceCount, len(placeCounts))
	for i, v := range placeCounts {
		result[i] = PlaceCount{Place: v.Place, Count: v.Count}
	}
	return result, nil
}

func getEmigrations(db DbSelecter, place, year string) ([]PlaceCount, error) {
	placeCounts := []EmigrationPlaceCount{}
	err := db.Select(&placeCounts, "SELECT place_to, count FROM emigrations WHERE place_from = ? and year = ?", place, year)
	if err != nil {
		return nil, err
	}

	result := make([]PlaceCount, len(placeCounts))
	for i, v := range placeCounts {
		result[i] = PlaceCount{Place: v.Place, Count: v.Count}
	}
	return result, nil
}

var re = regexp.MustCompile(",\"(\\d+)\"\\]")

func removeCountQuotes(s string) string {
	return re.ReplaceAllString(s, ",$1]")
}

func asResult(placeCounts []PlaceCount, filter string) [][]string {
	filterLevels := getLevels(filter)
	aggs := map[string]int{}
	for _, pc := range placeCounts {
		placeLevels := getLevels(clean(pc.Place))
		if endsWithCountry(placeLevels) && containsPlace(filterLevels, placeLevels) {
			place := constructPlace(placeLevels, len(filterLevels)+1)
			aggs[place] += pc.Count
		}
	}
	result := [][]string{}
	for k, v := range aggs {
		result = append(result, []string{k, strconv.Itoa(v)})
	}
	return result
}

func clean(place string) string {
	// if place ends with space-Territory, United States, remove space-Territory
	if strings.HasSuffix(strings.ToLower(place), " territory, united states") {
		place = place[0:len(place)-25] + place[len(place)-15:]
	}
	return place
}

func getLevels(place string) []string {
	levels := []string{}
	if len(place) > 0 {
		for _, level := range strings.Split(place, ",") {
			levels = append(levels, strings.TrimSpace(level))
		}
	}
	return levels
}

func endsWithCountry(levels []string) bool {
	return countryMap[strings.ToLower(levels[len(levels)-1])]
}

func containsPlace(super []string, sub []string) bool {
	extra := len(sub) - len(super)
	if extra <= 0 {
		return false
	}
	for i := range super {
		if strings.ToLower(super[i]) != strings.ToLower(sub[i+extra]) {
			return false
		}
	}
	return true
}

func constructPlace(levels []string, max int) string {
	return strings.Join(levels[len(levels)-max:], ", ")
}
