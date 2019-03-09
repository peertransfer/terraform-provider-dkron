package dkron

import (
	"log"
)

type Config struct {
	Host string
}

func (c *Config) Client() {
	// TODO implement client
	log.Print("[INFO] Dkron client configured")
}