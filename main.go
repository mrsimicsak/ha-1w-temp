package main

import (
	"fmt"
	"time"

	"github.com/yryz/ds18b20"

	"github.com/spf13/viper"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// CtoF converts celsius to fahrenheit
func CtoF(c float64) float64 {
	return (c * 1.8) + 32
}

func main() {

	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")      // optionally look for config in the working directory
	err := viper.ReadInConfig()   // Find and read the config file
	if err != nil {               // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	viper.SetDefault("MQTT.port", "1883")

	// Connect to the MQTT server
	opts := mqtt.NewClientOptions().AddBroker("tcp://" + viper.GetString("MQTT.broker") + ":" + viper.GetString("MQTT.port"))
	opts.SetClientID(viper.GetString("MQTT.clientid"))
	opts.SetKeepAlive(30 * time.Second)
	opts.SetPingTimeout(10 * time.Second)

	opts.SetUsername(viper.GetString("MQTT.username"))

	opts.SetPassword(viper.GetString("MQTT.password"))

	opts.SetWill("home/status/ds18b20-gateway-01", "dead", 1, true)

	c := mqtt.NewClient(opts)

	sensors, err := ds18b20.Sensors()
	if err != nil {
		panic(err)
	}

	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	c.Publish("home/status/ds18b20-gateway-01", 1, true, "running")

	fmt.Printf("sensor IDs: %v\n", sensors)

	for _, sensor := range sensors {

		topic := "homeassistant/sensor/ds18b20-gateway-01/" + sensor + "/config"

		config := `{
            "unit_of_measurement": "°F",
            "payload_not_available": "dead",
            "expire_after": 125,
            "force_update": true,
            "name": "Temp Sensor ` + sensor + `",
            "state_topic": "home/sensor/ds18b20-gateway-01/` + sensor + `",
            "availability_topic": "home/status/ds18b20-gateway-01",
            "unique_id": "ds18b20-gateway-01-` + sensor + `-f",
            "payload_available": "running"
            }`

		c.Publish(topic, 0, true, config)

	}

	for true {
		for _, sensor := range sensors {
			t, err := ds18b20.Temperature(sensor)
			if err == nil {
				fmt.Printf("sensor: %s temperature: %.2f°C, %.2f°F\n", sensor, t, CtoF(t))

				topic := "home/sensor/ds18b20-gateway-01/" + sensor
				c.Publish(topic, 0, false, fmt.Sprintf("%f", CtoF(t)))
			}
		}
		time.Sleep(1 * time.Minute)
	}
}
