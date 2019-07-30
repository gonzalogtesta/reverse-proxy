package routes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"time"
)

/*
LimitType allows to specify the type of limit to apply
*/
type LimitType string

/*
Constants for Limit types
*/
const (
	OriginIP                LimitType = "OriginIP"
	DestinationPath         LimitType = "Path"
	OriginIPDestinationPath LimitType = OriginIP + DestinationPath
	OriginIPUserAgent       LimitType = OriginIP + "UserAgent"
)

/*
RouteConfig allows to config
*/
type RouteConfig struct {
	Path      string `json:path`
	Server    string `json:server`
	Limit     int64  `json:limit`
	Time      time.Duration
	strTime   string    `json:time`
	LimitType LimitType `json:limitType`
}

/*
RouteConfig allows to config
*/
type dataRouteConfig struct {
	Path      string    `json:path`
	Server    string    `json:server`
	Limit     int64     `json:limit`
	Time      string    `json:time`
	LimitType LimitType `json:limitType`
}

/*
ReadFileRoute reads a JSON file and generates a RouteConfig
*/
func ReadFileRoute(filename string) RouteConfig {

	path, err := filepath.Abs(filename)
	if err != nil {
		log.Fatalf("Unable to find file: %s", err)
	}

	file, err := ioutil.ReadFile(path)

	if err != nil {
		log.Fatalf("Unable to read file: %s", err)
	}

	return parseRoute(file)
}

func parseRoute(file []byte) RouteConfig {
	data := dataRouteConfig{}

	err := json.Unmarshal([]byte(file), &data)

	if err != nil {
		log.Fatalf("Unable to parse JSON: %s", err)
	}

	duration, _ := time.ParseDuration(data.Time) // fix issue with unmarshal of durations

	fmt.Println("Parsed route: ", data)

	route := RouteConfig{
		Path:      data.Path,
		Server:    data.Server,
		Limit:     data.Limit,
		Time:      duration,
		LimitType: data.LimitType,
	}

	return route
}
