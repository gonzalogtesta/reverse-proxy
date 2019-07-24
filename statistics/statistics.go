package statistics

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

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

 */
func getIP(r *http.Request) (ip string) {
	ip = r.Header.Get("X-Forwarded-For")
	if ip == "" {
		if strings.Count(r.RemoteAddr, ":") < 2 {
			ip = strings.Split(r.RemoteAddr, ":")[0]
		} else {
			ip = strings.Split(r.RemoteAddr, "]")[0] + "]"
		}
	}
	return
}

/*
Track tracks
*/
func (m *Metrics) Track(r *http.Request) {
	fmt.Println("Ok:")
	var duration, _ = time.ParseDuration("5m")

	var keyname = "user_request:" + getIP(r)

	fmt.Println("Key name: " + keyname)
	_, havit := m.redisConn.Info(keyname)
	fmt.Println("Have it: " + keyname)
	if false && havit != nil {
		m.redisConn.CreateKey(keyname, duration)
		m.redisConn.CreateKey(keyname+"_avg", 0)
		m.redisConn.CreateRule(keyname, redis.AvgAggregation, 60, keyname+"_avg")
	}

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
Get gets
*/
func (m *Metrics) Get(duration time.Duration) (info map[string]interface{}, err error) {
	/*
		resp, err := m.redisConn.Info("user_request:" + getIP(r))
		fmt.Println(resp)
		if err != nil {
			// handle error
			fmt.Println("Fail: ", err)
		}
	*/
	info = make(map[string]interface{})

	keys, _ := m.redisConn.Keys()
	for _, key := range keys {
		data, _ := m.redisConn.AggRange(key, time.Now().Add(time.Minute*-30).UnixNano()/1e6, time.Now().UnixNano()/1e6, redis.CountAggregation, 1000)
		info[key] = data
	}

	return info, nil
}

/*
Get gets series
*/
func (m *Metrics) GetSerie(key string, duration time.Duration) (info [][]int64, err error) {
	/*
		resp, err := m.redisConn.Info("user_request:" + getIP(r))
		fmt.Println(resp)
		if err != nil {
			// handle error
			fmt.Println("Fail: ", err)
		}
	*/

	data, _ := m.redisConn.AggRange(key, time.Now().Add(time.Minute*-30).UnixNano()/1e6, time.Now().UnixNano()/1e6, redis.CountAggregation, 1000)
	for _, key := range data {
		info = append(info, []int64{key.Timestamp, int64(key.Value)})
	}

	return info, nil
}

/*
GetForPeriod gets
*/
func (m *Metrics) GetForPeriod(r *http.Request, duration time.Duration) (sum int64, err error) {

	var keyname = "user_request:" + getIP(r)
	resp, _ := m.redisConn.AggRange(keyname,
		time.Now().Add(time.Second*-30).UnixNano()/1e6,
		time.Now().UnixNano()/1e6, redis.CountAggregation,
		int(time.Millisecond*30))
	sum = 0
	for _, item := range resp {
		fmt.Println(item.Value)
		// val, _ := item.Value
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
SendMetrics allows to send metric information from a HTTP Request to a AMQP instance.

*/
func (m *Metrics) SendMetrics(r *http.Request) {
	/*
		fmt.Println(r.RemoteAddr)
		fmt.Println(r.Header.Get("X-Forwarded-For"))
		fmt.Println(r.Header.Get("User-Agent"))
		fmt.Println(r.Method)
		fmt.Println(r.RequestURI)
		fmt.Println(r.Proto)

		obj := map[string]string{
			"remoteAddr":      r.RemoteAddr,
			"x-forwarded-for": r.Header.Get("X-Forwarded-For"),
			"user-agent":      r.Header.Get("User-Agent"),
			"method":          r.Method,
			"requestURI":      r.RequestURI,
			"proto":           r.Proto,
		}

		//Publish to the queue

		jsonString, err := json.Marshal(obj)

		err = m.ch.Publish(
			"proxy.new-request", //exchange
			m.q.Name,            //routing key
			false,               //mandatory
			false,               //immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(jsonString),
			})

		failOnError(err, "Failed to publish a message ", "Published the message")
	*/
}
