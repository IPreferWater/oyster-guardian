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

	MQTT "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

const (
	TCP = "tcp"
	SSL = "ssl"
)

var (
	protocol     = os.Getenv("MQTT_PROTOCOL")
	host         = os.Getenv("MQTT_HOST")
	port         = os.Getenv("MQTT_PORT")
	clientID         = os.Getenv("MQTT_CLIENT_ID")
	username         = os.Getenv("MQTT_USERNAME")
	password         = os.Getenv("MQTT_PASSWORD")
)

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

		topic := ("oyster-guardian/random")
		if token := c.Subscribe(topic, byte(2), mqttRandom); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
	}

	client := MQTT.NewClient(connOpts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Errorf("error can't connect", token.Error())
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

func mqttRandom(client MQTT.Client, message MQTT.Message) {
	payload := string(message.Payload())
	log.Debug(payload)
}