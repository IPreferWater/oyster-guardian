package mqtt

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/IPreferWater/oyster-guardian/service"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

const (
	TCP            = "tcp"
	SSL            = "ssl"
	TOPIC_DETECTED = "oyster-guardian/detected"
	TOPIC_THREAT   = "oyster-guardian/threat"
)

var (
	/*protocol     = os.Getenv("MQTT_PROTOCOL")
	host         = os.Getenv("MQTT_HOST")
	port         = os.Getenv("MQTT_PORT")
	clientID         = os.Getenv("MQTT_CLIENT_ID")
	username         = os.Getenv("MQTT_USERNAME")
	password         = os.Getenv("MQTT_PASSWORD")*/
	protocol = "TCP"
	host     = "localhost"
	port     = "1883"
	clientID = "clientID"
	username = "username"
	password = "password"

	//MqttClient MQTT.Client
)

type MqttStream struct {
	client MQTT.Client
}

func InitMqtt() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	url := fmt.Sprintf("%s://%s:%s", protocol, host, port)
	connOpts := MQTT.NewClientOptions().AddBroker(url).SetClientID(clientID).SetCleanSession(false).SetAutoReconnect(true).SetConnectionLostHandler(connectionLostHandler)
	connOpts.SetUsername(username)
	connOpts.SetPassword(password)

	tlsConfig := getTlsConfig(protocol)
	connOpts.SetTLSConfig(&tlsConfig)

	connOpts.OnConnect = func(c MQTT.Client) {
		log.Infof("Connected to %s\n", url)

		if token := c.Subscribe(TOPIC_DETECTED, byte(2), subscribeTopicDetected); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}

		if token := c.Subscribe(TOPIC_THREAT, byte(2), subscribeTopicThreat); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
	}

	client := MQTT.NewClient(connOpts)
	service.Stream = MqttStream{
		client: client,
	}
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Errorf("can't connect MQTT", token.Error())
		time.Sleep(5 * time.Second)
		InitMqtt()
	}

	<-c
}

func getTlsConfig(protocol string) tls.Config {

	if strings.EqualFold(SSL, protocol) {
		return getSSLConfig()
	}

	//if it's not ssl, we set the default value to tcp style
	return getTCPConfig()
}

func getTCPConfig() tls.Config {
	return tls.Config{
		ClientCAs:          nil,
		InsecureSkipVerify: true,
	}
}

func getSSLConfig() tls.Config {
	certpool := x509.NewCertPool()
	pemCerts, err := ioutil.ReadFile("/etc/oyster-guardian/ca.crt")
	if err == nil {
		certpool.AppendCertsFromPEM(pemCerts)
	}

	cert, err := tls.LoadX509KeyPair("/etc/oyster-guardian/server.crt", "/etc/oyster-guardian/server.key")
	if err != nil {
		panic(err)
	}

	return tls.Config{
		RootCAs:            certpool,
		ClientCAs:          nil,
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
		ClientAuth:         tls.NoClientCert}
}

func connectionLostHandler(client MQTT.Client, err error) {
	log.Errorf("Connection lost, reason: %v", err)
}

func subscribeTopicDetected(client MQTT.Client, message MQTT.Message) {
	payload := string(message.Payload())
	log.Infof("topic %s consume message %s", TOPIC_DETECTED, payload)
	service.HandleTopicDetected(payload)
}

func subscribeTopicThreat(client MQTT.Client, message MQTT.Message) {
	payload := string(message.Payload())
	log.Infof("topic %s consume message %s", TOPIC_THREAT, payload)
	service.HandleTopicThreat(payload)
}

func (p MqttStream) PublishTopicDetected(payload string) error {
	token := p.client.Publish(TOPIC_DETECTED, byte(2), true, payload)
	return token.Error()
}

func (p MqttStream) PublishTopicThreat(payload string) error {
	token := p.client.Publish(TOPIC_THREAT, byte(2), true, payload)
	return token.Error()
}
