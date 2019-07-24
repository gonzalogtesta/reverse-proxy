package redis

import (
	"time"

	redistimeseries "github.com/RedisTimeSeries/redistimeseries-go"
	redigo "github.com/gomodule/redigo/redis"
)

/*
Client extends Redist Timeseries client, to provide additionals operations
*/
type Client struct {
	redistimeseries.Client
	Pool redistimeseries.ConnPool
	Name string
}

/*
Aggregation constants
*/
const (
	AvgAggregation   redistimeseries.AggregationType = redistimeseries.AvgAggregation
	SumAggregation   redistimeseries.AggregationType = redistimeseries.SumAggregation
	MinAggregation   redistimeseries.AggregationType = redistimeseries.MinAggregation
	MaxAggregation   redistimeseries.AggregationType = redistimeseries.MaxAggregation
	CountAggregation redistimeseries.AggregationType = redistimeseries.CountAggregation
	FirstAggregation redistimeseries.AggregationType = redistimeseries.FirstAggregation
	LastAggregation  redistimeseries.AggregationType = redistimeseries.LastAggregation
)

/*
NewClient returns a new client
*/
func NewClient(addr, name string, authPass *string) *Client {
	cli := redistimeseries.NewClient(addr, name, authPass)
	return &Client{
		Pool:   cli.Pool,
		Name:   cli.Name,
		Client: *cli,
	}
}

/*
IncBy increments by one the given key
*/
func (client *Client) IncBy(key string, value int64, reset time.Duration) (reponse string, err error) {
	conn := client.Pool.Get()
	defer conn.Close()
	// TS.INCRBY key [VALUE] [RESET] [TIME_BUCKET]
	return redigo.String(conn.Do("TS.INCRBY", key, value, "RESET", 1000, "RETENTION", 1000))
}

/*
Keys returns all the keys for the pattern given
*/
func (client *Client) Keys() (keys []string,
	err error) {
	conn := client.Pool.Get()
	defer conn.Close()
	info, err := conn.Do("KEYS", "*")

	keys, err = redigo.Strings(info, err)
	return keys, err
}
