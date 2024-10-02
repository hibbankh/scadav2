package mqtt

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type (
	ClientStruct struct {
		mq1     mqtt.Client
		mq2     mqtt.Client
		servers []string
	}
)

var mqttBroker mqtt.Client

func NewMQTTClient(server ...string) *ClientStruct {

	uname := os.Getenv("MQ_USER")
	passwd := os.Getenv("MQ_PASS")

	// tls := NewTLSConfig()
	rand.Seed(time.Now().UnixNano())
	log.Printf("MQTT Server: %v\n", server)

	client1 := mqtt.NewClientOptions().AddBroker(server[0]).
		SetClientID(fmt.Sprintf("vecto-%d", rand.Int())).
		SetUsername(uname).
		SetPassword(passwd).
		// SetTLSConfig(tls).
		// SetDefaultPublishHandler(MsgHandlerBrokerNBIOT).
		SetAutoReconnect(true).
		SetCleanSession(false)

	a := mqtt.NewClient(client1)

	if token := a.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	mqttBroker = a
	return &ClientStruct{
		mq1: a,
	}
}

func (c *ClientStruct) GetClient1() mqtt.Client {
	return c.mq1
}

func GetMqttClient() mqtt.Client {
	return mqttBroker
}

func MsgHandlerInstReadingLog(client mqtt.Client, msg mqtt.Message) {
	payload := msg.Payload()
	fmt.Println("received")
	InstReadingLog(payload, msg.Topic())
}

func MsgHandlerAlarmLog(client mqtt.Client, msg mqtt.Message) {
	payload := msg.Payload()
	HandleAlarmLog(payload, msg.Topic())
}
