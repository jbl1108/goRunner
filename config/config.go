package config

import (
	"log"
	"strings"

	"github.com/magiconair/properties"
)

type Config struct {
	prop *properties.Properties
}

const CONFIG_FILE = "config.properties"

func NewConfig() *Config {
	c := new(Config)
	var err error
	configPath := "/gosecret/config/" + CONFIG_FILE
	log.Printf("Loading config file from %s", configPath)

	c.prop, err = properties.LoadFile(configPath, properties.UTF8)
	if err != nil {
		log.Printf("Failed to load config file: %v. Using default values.", err)
		c.prop = properties.NewProperties()
	}
	log.Printf("Using the following configuration:")
	for _, key := range c.prop.Keys() {
		value, _ := c.prop.Get(key)
		log.Printf("%s = %s", key, value)
	}

	return c
}

func sanitize(s string) string {
	// trim spaces and tabs
	return strings.TrimSpace(s)
}

func (config *Config) MQTTAddress() string {
	return sanitize(config.prop.GetString("mqtt_address", "localhost:1883"))
}
func (config *Config) MQTTUsername() string {
	return sanitize(config.prop.GetString("mqtt_username", "mqtt-user"))
}
func (config *Config) MQTTPassword() string {
	return sanitize(config.prop.GetString("mqtt_password", "mqtt-password"))
}

func (config *Config) KeyValueDBURL() string {
	return sanitize(config.prop.GetString("gokeyvaluestore_url", "http://localhost:9091"))
}
func (config *Config) RestAddress() string {
	return sanitize(config.prop.GetString("own_rest_address", ":9092"))
}
