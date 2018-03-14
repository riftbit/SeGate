package main

import (
	"crypto/rsa"

	"database/sql"

	"github.com/Sirupsen/logrus"
)

var (
	version, build, buildDate string
)

const ServerName = "SecGW by ErgoZ"
const PoweredBy = "Gcode_"

type conf struct {
	System struct {
		MaxThreads int    `yaml:"maxThreads"`
		ListenOn   string `yaml:"listenOn"`
		AESKey     string `yaml:"aesKey"`
	}

	Nodes []struct {
		NodeName string `yaml:"name"`
		NodeURL  string `yaml:"url"`
		AESKey   string `yaml:"aesKey"`
	}

	Clickhouse struct {
		ConnString string `yaml:"connString"`
		DBName     string `yaml:"dbName"`
		IsEnabled  bool   `yaml:"isEnabled"`
	}

	Log struct {
		Formatter       string `yaml:"formatter"` //text, json
		LogLevel        string `yaml:"logLevel"`  // panic, fatal, error, warn, warning, info, debug
		DisableColors   bool   `yaml:"disableColors"`
		TimestampFormat string `yaml:"timestampFormat"`
	}
}

type nodeElement struct {
	NodeName string
	NodeURL  string
	AESKey   []byte
}

var NodesList map[string]nodeElement

//Logger ...
var Logger *logrus.Logger

//Configs ...
var config conf

var configPath string

var PublicKey *rsa.PublicKey

var clickHouseDB *sql.DB

//var PrivateKey *rsa.PrivateKey
