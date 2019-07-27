package redis

import (
	"fmt"

	redistimeseries "github.com/RedisTimeSeries/redistimeseries-go"
	"github.com/gomodule/redigo/redis"
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
func (client *Client) IncBy(key string, value int64) (reponse string, err error) {
	conn := client.Pool.Get()
	defer conn.Close()
	// TS.INCRBY key [VALUE] [RESET] [TIME_BUCKET]
	return redigo.String(conn.Do("TS.INCRBY", key, value)) //, "RESET", 0, "RETENTION", 1000))
}

/*
Keys returns all the keys for the pattern given
*/
func (client *Client) Keys() (keys []string, err error) {
	conn := client.Pool.Get()
	defer conn.Close()
	info, err := conn.Do("KEYS", "*")

	keys, err = redigo.Strings(info, err)
	return keys, err
}

/*
KeysNames returns all the keys for the pattern given
*/
func (client *Client) KeysNames(pattern string) (keys []string, err error) {
	conn := client.Pool.Get()
	defer conn.Close()
	n := 0
	for {
		fmt.Println("Pattern: ", pattern)
		arr, err := redis.Values(conn.Do("SCAN", n, "MATCH", pattern))
		if err != nil {
			return keys, fmt.Errorf("error retrieving '%s' keys", pattern)
		}
		fmt.Println("Arr: ", arr)
		n, _ = redis.Int(arr[0], nil)
		k, _ := redis.Strings(arr[1], nil)
		keys = append(keys, k...)
		if n == 0 {
			break
		}
	}

	return keys, err
}

func (client *Client) TrackTime(keyname string, value float64) {
	key := fmt.Sprintf("ls:%s", keyname)
	conn := client.Pool.Get()
	defer conn.Close()
	conn.Do("RPUSH", key, value)
}

func (client *Client) LLen(keyname string) string {
	conn := client.Pool.Get()
	defer conn.Close()
	len, _ := redis.String(conn.Do("LLEN", keyname))
	return len
}
