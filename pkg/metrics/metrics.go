package metrics

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"meli-proxy/pkg/keys"
	"meli-proxy/pkg/routes"
	"meli-proxy/utils/redis"
)

/*
Metrics represents a metrics handler. It allows to track all information.

*/
type Metrics struct {
	ctx       context.Context
	redisConn *redis.Client
}

func failOnError(err error, msgerr string) {
	if err != nil {
		log.Fatalf("%s: %s", msgerr, err)
	}
}

/*
Track tracks
*/
func (m *Metrics) Track(route routes.RouteConfig, r *http.Request) {

	var keyname = keys.GenerateKey(route, r)

	now := time.Now().UnixNano() / 1e6 // now in ms
	_, err := m.redisConn.Add(keyname, now, 1)

	if err != nil {
		fmt.Println("Error:", err)
	}
}

/*
CreateKey creates key with duration
*/
func (m *Metrics) CreateKey(key string, duration string) {
	var timeDuration, _ = time.ParseDuration(duration)
	_, havit := m.redisConn.Info(key)
	fmt.Println("Have it: " + key)
	if havit != nil {
		m.redisConn.CreateKey(key, timeDuration)
	}
}

/*
Hit allows to track hits to the server
*/
func (m *Metrics) Hit(r *http.Request) {
	var keyname = "hits"
	now := time.Now().UnixNano() // 1e6 // now in ms
	m.redisConn.IncBy(keyname, now)
}

/*
SendCode allows to track status code of the server
*/
func (m *Metrics) SendCode(code int, startTime time.Time) {
	var keyname = fmt.Sprintf("response_%d", code)
	now := time.Now().UnixNano()
	m.redisConn.IncBy(keyname, now)

	seconds := time.Now().UnixNano() / int64(time.Second)
	newkey := fmt.Sprintf("%s:time", keyname)
	_, err := m.redisConn.IncBy(newkey, seconds)
	if err != nil {
		fmt.Println("Error:", err)
	}
	m.redisConn.TrackTime(fmt.Sprintf("%s:%d", keyname, seconds), float64(time.Since(startTime))/float64(time.Millisecond))

}

func (m *Metrics) GetPercentile(percentile int, duration time.Duration) (info [][]int64, err error) {
	// TODO: calculate percentiles
	// Get all keys in a range (duration) e.g. using key name response_200:time
	// For each timestamp get size of list response_200:timestamp and calculate percentile
	return nil, nil
}

/*
Get gets
*/
func (m *Metrics) Get(duration time.Duration) (info map[string]interface{}, err error) {

	info = make(map[string]interface{})
	now := time.Now()
	fromTimestamp := now.Add(-duration).UnixNano() / 1e6
	toTimestamp := time.Now().UnixNano() / 1e6

	keys, _ := m.redisConn.Keys()
	for _, key := range keys {
		data, _ := m.redisConn.AggRange(key, fromTimestamp, toTimestamp, redis.CountAggregation, 1000)
		info[key] = data
	}

	return info, nil
}

/*
GetSerie gets series
*/
func (m *Metrics) GetSerie(key string, duration time.Duration) (info [][]int64, err error) {

	now := time.Now()
	fromTimestamp := now.Add(-duration).UnixNano() / 1e6
	toTimestamp := time.Now().UnixNano() / 1e6
	data, err := m.redisConn.AggRange(key, fromTimestamp, toTimestamp, redis.CountAggregation, 1000)
	if err != nil {
		return info, err
	}

	for _, key := range data {
		info = append(info, []int64{key.Timestamp, int64(key.Value)})
	}

	return info, nil
}

/*
GetForPeriod gets number of counts for a period
*/
func (m *Metrics) GetForPeriod(keyname string, duration time.Duration) (sum int64, err error) {

	now := time.Now()
	fromTimestamp := now.Add(-duration).UnixNano() / 1e6
	toTimestamp := time.Now().UnixNano() / 1e6
	resp, _ := m.redisConn.AggRange(keyname, fromTimestamp, toTimestamp, redis.CountAggregation,
		int(duration))
	sum = 0
	for _, item := range resp {
		sum += int64(item.Value)
	}

	return sum, nil
}

/*
Connect connects
*/
func (*Metrics) Connect(ctx context.Context) Metrics {

	m := Metrics{
		ctx:       ctx,
		redisConn: redis.NewClient("localhost:6379", "nohelp", nil),
	}

	return m
}

/*
NewMetrics connects
*/
func NewMetrics(ctx context.Context) Metrics {

	m := Metrics{
		ctx:       ctx,
		redisConn: redis.NewClient("localhost:6379", "nohelp", nil),
	}

	return m
}

/*
GetKeys get keys
*/
func (m *Metrics) GetKeys(generics []string) (keys []string) {

	for _, pattern := range generics {
		retrievedKeys, _ := m.redisConn.KeysNames(pattern)
		fmt.Println("Keys: ", retrievedKeys)
		for _, key := range retrievedKeys {
			keys = append(keys, key)
		}
	}

	return keys
}
