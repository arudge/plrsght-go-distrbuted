package main

import (
	"flag"
	"time"
	"strconv"
	"math/rand"
	"log"
	"git.target.com/plrsght-go-distrbuted/dto"
	"encoding/gob"
	"bytes"
	"git.target.com/plrsght-go-distrbuted/qutils"
	"github.com/streadway/amqp"
)

var url = "amqp://guest:guest@localhost:5672"

var name = flag.String("name", "sensor", "name of the sensor")
var freq = flag.Uint("freq", 5, "update frequency in cycles/sec")
var max = flag.Float64("max", 5., "maximum value for generated readings")
var min = flag.Float64("min", 1., "minimum value for generated readings")
var stepSize = flag.Float64("step", 0.1, "maxium allowable change per measurement")

var r = rand.New(rand.NewSource(time.Now().UnixNano()))
var nom = (*max - *min) / 2 + *min

var value = r.Float64() * (*max - *min) + *min

func main() {
	flag.Parse()

	conn, ch := qutils.GetChannel(url)
	defer ch.Close()
	defer conn.Close()

	dataQueue := qutils.GetQueue(*name, ch, false)

	publishQueueName(ch)

	discoveryQueue := qutils.GetQueue("", ch, true)
	ch.QueueBind(
		discoveryQueue.Name,
		"",
		qutils.SensorDiscoveryExchange,
		false,
		nil)

	go listenForDiscoveryRequests(discoveryQueue.Name, ch)

	dur, _ := time.ParseDuration(strconv.Itoa(1000 / int(*freq)) + "ms")

	signal := time.Tick(dur)

	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)

	for range signal {
		calcValue()
		reading := dto.SensorMessage{
			Name: *name,
			Value:value,
			Timestamp: time.Now(),
		}

		buf.Reset()
		enc = gob.NewEncoder(buf)
		enc.Encode(reading)

		msg := amqp.Publishing{
			Body: buf.Bytes(),
		}

		ch.Publish(
			"",
			dataQueue.Name,
			false,
			false,
			msg)

		log.Printf("Reading sent. Value: %v\n", value)
	}
}

func listenForDiscoveryRequests(name string, ch *amqp.Channel) {
	msgs, _ := ch.Consume(name, "", true, false, false, false, nil)

	for range msgs {
		publishQueueName(ch)
	}
}

func publishQueueName(ch *amqp.Channel) {
	msg := amqp.Publishing{Body: []byte(*name)}
	ch.Publish(
		"amq.fanout",
		"",
		false,
		false,
		msg)
}

func calcValue() {
	var maxStep, minStep float64

	if value < nom {
		maxStep = *stepSize
		minStep = -1 * *stepSize * (value - *min) / (nom - *min)
	} else {
		maxStep = *stepSize * (*max - value) / (*max - nom)
		minStep = -1 * *stepSize
	}

	value += r.Float64() * (maxStep - minStep) + minStep
}
