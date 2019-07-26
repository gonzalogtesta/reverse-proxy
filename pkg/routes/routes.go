package routes

import "time"

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
	UserAgent               LimitType = "UserAgent"
)

/*
RouteConfig allows to config
*/
type RouteConfig struct {
	Path      string
	Server    string
	Limit     int64
	Time      time.Duration
	LimitType LimitType
}
