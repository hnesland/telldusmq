package main

import (
	"bytes"
	"io"
	"log"
	"strings"
	"text/template"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

/*
#cgo LDFLAGS: -L. -ltelldus-core
#include <stdio.h>
#include <telldus-core.h>

void rawTelldusEvent(char *);

static inline void rawEvent(const char *data, int controllerId, int callbackId, void *context) {
	rawTelldusEvent((char*)data);
}

static inline void initTelldus() {
	tdRegisterRawDeviceEvent(&rawEvent, NULL);
}
*/
import "C"

// TelldusEvent describes telldus events ..
type TelldusEvent struct {
	Class    string
	Protocol string
	Model    string
	Code     string
	House    string
	Unit     string
	Group    string
	Method   string
	Id       string
	Temp     string
	Humidity string
	Value    string
	DataType string
}

var mqttClient MQTT.Client

//export rawTelldusEvent
func rawTelldusEvent(str *C.char) {
	data := strings.Split(C.GoString(str), ";")
	event := &TelldusEvent{
		Id:       "0",
		House:    "0",
		Unit:     "0",
		Code:     "0",
		Group:    "0",
		Temp:     "0",
		Humidity: "0",
		Method:   "0",
		Value:    "0",
		DataType: ""}
	for _, elm := range data {
		if len(elm) != 0 {
			propval := strings.Split(elm, ":")
			switch propval[0] {
			case "class":
				event.Class = propval[1]
				break
			case "protocol":
				event.Protocol = propval[1]
				break
			case "model":
				event.Model = propval[1]
				break
			case "code":
				event.Code = propval[1]
				break
			case "house":
				event.House = propval[1]
				break
			case "unit":
				event.Unit = propval[1]
				break
			case "group":
				event.Group = propval[1]
				break
			case "method":
				event.Method = propval[1]
				break
			case "id":
				event.Id = propval[1]
				break
			case "temp":
				event.Temp = propval[1]
				break
			case "humidity":
				event.Humidity = propval[1]
				break
			}
		} else {
			var topicTemplate string
			var payloadTemplate string
			if event.Class == "command" {
				topicTemplate = viper.GetString("Mqtt.Events.PublishTopic")
				payloadTemplate = viper.GetString("Mqtt.Events.PublishPayload")

				turnOn := viper.GetString("Tellstick.MapTurnOnTo")
				turnOff := viper.GetString("Tellstick.MapTurnOffTo")

				if len(turnOn) > 0 && event.Method == "turnon" {
					event.Method = turnOn
				}

				if len(turnOff) > 0 && event.Method == "turnoff" {
					event.Method = turnOff
				}
			} else {
				topicTemplate = viper.GetString("Mqtt.Sensors.PublishTopic")
				payloadTemplate = viper.GetString("Mqtt.Sensors.PublishPayload")
				event.Value = event.Temp
				event.DataType = "temp"
			}

			var topicString string
			var payloadString string

			topicString = parseTemplate(topicTemplate, event)
			payloadString = parseTemplate(payloadTemplate, event)

			var token MQTT.Token

			log.Printf("Publish to '%s' with '%s'\n", topicString, payloadString)
			token = mqttClient.Publish(topicString, 0, false, payloadString)
			token.Wait()

			// Send a duplicate event for humidity
			if viper.GetBool("Tellstick.SplitTemperatureAndHumidity") && event.Class == "sensor" {
				event.DataType = "humidity"
				event.Value = event.Humidity
				topicString = parseTemplate(topicTemplate, event)
				payloadString = parseTemplate(payloadTemplate, event)

				log.Printf("Publish to '%s' with '%s'\n", topicString, payloadString)
				token = mqttClient.Publish(topicString, 0, false, payloadString)
				token.Wait()
			}
		}
	}
}

func parseTemplate(templateString string, event *TelldusEvent) string {
	tmpl, err := template.New("template").Parse(templateString)
	if err != nil {
		log.Panicf("Error parsing template: %v", err)
	}

	var tmplBuffer bytes.Buffer
	tmplWriter := io.Writer(&tmplBuffer)
	err = tmpl.Execute(tmplWriter, event)

	if err != nil {
		log.Panicf("Error executing template: %v", err)
	}

	return tmplBuffer.String()
}

func setupMqtt() {
	opts := MQTT.NewClientOptions()

	opts.AddBroker(viper.GetString("Mqtt.Broker"))
	opts.SetClientID(viper.GetString("Mqtt.ClientId"))
	opts.SetUsername(viper.GetString("Mqtt.Username"))
	opts.SetPassword(viper.GetString("Mqtt.Password"))

	mqttClient = MQTT.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		log.Panicf("Unable to connect to MQTT: %v", token.Error())
	}
}

func setupConfiguration() {
	viper.SetConfigName("telldusmq")
	viper.AddConfigPath("/etc/telldusmq/")
	viper.AddConfigPath("$HOME/.telldusmq/")
	viper.AddConfigPath("./")

	configError := viper.ReadInConfig()
	if configError != nil {
		log.Panicf("Error reading configuration: %v\n", configError)
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Println("Reloading configuration", e.Name)
		// TODO: Reconnect mqtt broker if connection params changes?
	})
}

func main() {
	log.Println("Started Message Queue for Telldus Core")
	setupConfiguration()
	setupMqtt()
	C.initTelldus()

	for {
		time.Sleep(30 * time.Second)
	}
}
