package config

import (
	"log"
	"os"
	"strings"

	"github.com/magiconair/properties"
)

type Config struct {
	prop *properties.Properties
}

const CONFIG_FILE = "config.conf"

func NewConfig() *Config {
	c := new(Config)

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	configPath := dir + "/" + CONFIG_FILE
	log.Printf("Loading config file from %s", configPath)
	c.prop, err = properties.LoadFile(configPath, properties.UTF8)
	if err != nil {
		log.Printf("Failed to load config file: %v. Using default values.", err)
		c.prop = properties.NewProperties()
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
	return sanitize(config.prop.GetString("keyvaluedb_url", "http://localhost:8086"))
}
func (config *Config) RestAddress() string {
	return sanitize(config.prop.GetString("rest_address", ":8081"))
}
