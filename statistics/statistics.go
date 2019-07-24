package statistics

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	redis "github.com/RedisTimeSeries/redistimeseries-go"
	"github.com/streadway/amqp"
)

/*
Metrics represents a metrics handler. It allows to track all information.

*/
type Metrics struct {
	ctx       context.Context
	conn      *amqp.Connection
	ch        *amqp.Channel
	q         amqp.Queue
	redisConn *redis.Client
}

func failOnError(err error, msgerr string, msgsuc string) {
	if err != nil {
		log.Fatalf("%s: %s", msgerr, err)

	} else {
		fmt.Printf("%s\n", msgsuc)
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

	var keyname = "user_request:" + getIP(r) //+ ":" + r.RequestURI //+ "route:" + r.RequestURI

	fmt.Println("Key name: " + keyname)
	_, havit := m.redisConn.Info(keyname)
	fmt.Println("Have it: " + keyname)
	if false && havit == nil {
		m.redisConn.CreateKey(keyname, duration)
		m.redisConn.CreateKey(keyname+"_avg", 0)
		m.redisConn.CreateRule(keyname, redis.AvgAggregation, 60, keyname+"_avg")
	}
	labels := map[string]string{}
	//labels["ip"] = getIP(r)
	//labels["path"] = r.RequestURI
	now := time.Now().UnixNano() / 1e6 // now in ms
	_, err := m.redisConn.Add(keyname, now, 1, labels)
	//_, err := m.redisConn.IncBy(keyname, 1, time.Second)
	if err != nil {
		fmt.Println("Error:", err)
	}
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
		data, _ := m.redisConn.AggRange(key, time.Now().Add(time.Second*-30).UnixNano()/1e6, time.Now().UnixNano()/1e6, redis.CountAggregation, int(time.Millisecond*30))
		info[key] = data
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

	fmt.Println("Connecting to RabbitMQ ...")
	conn, err := amqp.Dial("") //Insert the  connection string
	failOnError(err, "RabbitMQ connection failure", "RabbitMQ Connection Established")

	//Connect to the channel

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel", "Opened the channel")

	//Declare the queue where messages need to be sent. Queue will be created if not already there
	q, err := ch.QueueDeclare(
		"Proxy", //name
		true,    //durable
		false,   //delete when unused
		false,   //exclusive
		false,   //no-wait
		nil,     //arguements
	)

	failOnError(err, "Failed to declare the queue", "Declared the queue")

	m := Metrics{
		ctx:  ctx,
		conn: conn,
		ch:   ch,
		q:    q,
	}
	go func() {
		<-ctx.Done()
		defer m.conn.Close()
		defer m.ch.Close()
	}()

	m.redisConn = redis.NewClient("localhost:6379", "nohelp", nil)

	return m
}

/*
SendMetrics allows to send metric information from a HTTP Request to a AMQP instance.

*/
func (m *Metrics) SendMetrics(r *http.Request) {

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

}
