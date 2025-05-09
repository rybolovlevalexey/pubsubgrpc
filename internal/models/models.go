package models

import (
	"log"
)

type Config struct{
	PubSubgRPCPort int
}


type ServerSettings struct{
	Cfg *Config
	Log *log.Logger
}