package metrics

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sort"
	"sync"
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
	//fmt.Println("New key: ", newkey)
	_, err := m.redisConn.IncBy(newkey, seconds)
	if err != nil {
		fmt.Println("Error:", err)
	}
	m.redisConn.TrackTime(fmt.Sprintf("%s:%d", keyname, seconds), float64(time.Since(startTime))/float64(time.Millisecond))

}

func (m *Metrics) getPercentileValue(percentile int, timestamp int64, value float64) []float64 { //, out chan []float64, wg *sync.WaitGroup) {
	index := float64(percentile) / 100.0 * value

	items := []float64{}
	keyname := fmt.Sprintf("ls:%s:%d", "response_200", timestamp/1000)
	if index == float64(int64(index)) {
		actual := int64(index)
		prev := actual - 1

		response1, _ := m.redisConn.GetFromSortedList(keyname, prev-1, 1)
		if len(response1) == 0 {
			return nil
		}
		response2, _ := m.redisConn.GetFromSortedList(keyname, actual-1, 1)
		items = []float64{
			float64(timestamp),
			float64((response1[0] + response2[0]) / 2),
		}
	} else {
		round := int64(index)
		response, _ := m.redisConn.GetFromSortedList(keyname, round, 1)
		if len(response) == 0 {
			return nil
		}
		items = []float64{
			float64(timestamp),
			float64(response[0]),
		}

	}
	// fmt.Println(items)
	return items
}

type Job struct {
	Percentile int
	Timestamp  int64
	Value      float64
}

func (m *Metrics) Calculate(id int, jobs <-chan Job, results chan<- []float64, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		result := m.getPercentileValue(job.Percentile, job.Timestamp, job.Value)
		if result != nil {
			results <- result
		}
	}
}

func (m *Metrics) GetPercentile(percentile int, duration time.Duration) (info [][]float64, err error) {
	// TODO: calculate percentiles for any key
	// Get all keys in a range (duration) e.g. using key name response_200:time
	// For each timestamp get size of list response_200:timestamp and calculate percentile
	now := time.Now()
	fromTimestamp := now.Add(-duration).UnixNano() / 1e6
	toTimestamp := time.Now().UnixNano() / 1e6
	data, _ := m.redisConn.AggRange("response_200:time", fromTimestamp, toTimestamp, redis.CountAggregation, 1000)

	var wg sync.WaitGroup

	jobs := make(chan Job, len(data))
	results := make(chan []float64, len(data))

	go func() {
		for _, c := range data {
			jobs <- Job{Percentile: percentile, Timestamp: c.Timestamp, Value: c.Value}
		}
		close(jobs)
	}()

	for i := 0; i < 15; i++ { // 15 consumers
		wg.Add(1)
		go m.Calculate(i, jobs, results, &wg)
	}

	wg.Wait()
	close(results)

	for val := range results {
		info = append(info, val)
	}

	sort.Slice(info, func(i, j int) bool {
		return info[i][0] > info[j][0]
	})

	return info, nil
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
