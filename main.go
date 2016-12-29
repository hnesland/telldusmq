package main

import (
	"bytes"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"io"
	"log"
	"strings"
	"text/template"
	"time"
)

/*
#cgo LDFLAGS: -L. -ltelldus-core
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <time.h>
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
		Method:   "0"}
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
			var publishTemplateString string
			var payloadTemplateString string
			if event.Class == "command" {
				publishTemplateString = viper.GetString("Mqtt.Events.PublishTopic")
				payloadTemplateString = viper.GetString("Mqtt.Events.PublishPayload")
			} else {
				publishTemplateString = viper.GetString("Mqtt.Sensors.PublishTopic")
				payloadTemplateString = viper.GetString("Mqtt.Sensors.PublishPayload")
			}

			publishTemplate, err := template.New("publish").Parse(publishTemplateString)
			if err != nil {
				log.Panicf("Error parsing publish template: %v", err)
			}

			var publishBuffer bytes.Buffer
			publishWriter := io.Writer(&publishBuffer)
			err = publishTemplate.Execute(publishWriter, event)

			if err != nil {
				log.Panicf("Error executing publish template: %v", err)
			}

			payloadTemplate, err := template.New("publish").Parse(payloadTemplateString)
			if err != nil {
				log.Panicf("Error parsing payload template: %v", err)
			}

			var payloadBuffer bytes.Buffer
			payloadWriter := io.Writer(&payloadBuffer)
			err = payloadTemplate.Execute(payloadWriter, event)

			if err != nil {
				log.Panicf("Error executing payload template: %v", err)
			}

			log.Printf("Publish to '%s' with '%s'\n", publishBuffer.String(), payloadBuffer.String())
			token := mqttClient.Publish(publishBuffer.String(), 0, false, payloadBuffer.String())
			token.Wait()
		}
	}
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

	viper.SetDefault("Mqtt.Broker", "tcp://localhost:1883")
	viper.SetDefault("Mqtt.ClientId", "TelldusMq")
	viper.SetDefault("Mqtt.Username", "")
	viper.SetDefault("Mqtt.Password", "")
	viper.SetDefault("Mqtt.Events.PublishTopic", "tellstick/events/{{.Protocol}}/{{.Model}}/{{.House}}/{{.Unit}}/{{.Group}}")
	viper.SetDefault("Mqtt.Events.PublishPayload", "{{.Method}}")
	viper.SetDefault("Mqtt.Events.SubscribeTopic", "tellstick/events")
	viper.SetDefault("Mqtt.Sensors.PublishTopic", "tellstick/sensors/{{.Protocol}}/{{.Model}}/{{.House}}/{{.Unit}}/{{.Group}}")
	viper.SetDefault("Mqtt.Sensors.PublishPayload", "Temp: {{.Temp}} Humidity: {{.Humidity}}")
	viper.SetDefault("Mqtt.Sensors.SubscribeTopic", "tellstick/sensors")

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
