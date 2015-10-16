package main

import (
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
)

type Params struct {
	Port string `envconfig:"PORT" required:"true"`
	DSN  string `required:"true"`
}

func main() {
	// parse params
	var params Params
	if err := envconfig.Process("fsmigrations", &params); err != nil {
		log.Fatal("Error parsing parameters", err)
	}
	log.Printf("Params %#v\n", params)

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

		placeCounts, err := getImmigrations(db, place, year)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, placeCounts)
	})

	e.Get("/emigrations", func(c *echo.Context) error {
		place := c.Query("place")
		year := c.Query("year")

		placeCounts, err := getEmigrations(db, place, year)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, placeCounts)
	})

	e.Run(":" + params.Port)
}

type DbSelecter interface {
	Select(dest interface{}, query string, args ...interface{}) error
}

type ImmigrationPlaceCount struct {
	Place string `json:"place" db:"place_from"`
	Count int    `json:"count" db:"count"`
}

func getImmigrations(db DbSelecter, place, year string) ([]ImmigrationPlaceCount, error) {
	placeCounts := []ImmigrationPlaceCount{}
	err := db.Select(&placeCounts, "SELECT place_from, count FROM immigrations WHERE place_to = ? and year = ?", place, year)
	if err != nil {
		return nil, err
	}
	return placeCounts, nil
}

type EmigrationPlaceCount struct {
	Place string `json:"place" db:"place_to"`
	Count int    `json:"count" db:"count"`
}

func getEmigrations(db DbSelecter, place, year string) ([]EmigrationPlaceCount, error) {
	placeCounts := []EmigrationPlaceCount{}
	err := db.Select(&placeCounts, "SELECT place_to, count FROM emigrations WHERE place_from = ? and year = ?", place, year)
	if err != nil {
		return nil, err
	}
	return placeCounts, nil
}
