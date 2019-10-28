package simpleflux

import (
	"log"
	"net/http"
	"strings"

	"github.com/spf13/viper"
)

// SimpleInfluxDB is a really primitive influxdb client
type SimpleInfluxDB struct {
	c        *http.Client
	db       string
	base     string
	username string
	password string
}

// NewSimpleInfluxDB gets a new simple influxdb client
func NewSimpleInfluxDB() *SimpleInfluxDB {
	return &SimpleInfluxDB{
		c:        &http.Client{},
		db:       viper.GetString("influxdbdatabase"),
		base:     strings.TrimSuffix(viper.GetString("influxdbhost"), `/`),
		username: viper.GetString("influxdbusername"),
		password: viper.GetString("influxdbpassword"),
	}
}

func (s *SimpleInfluxDB) Write(line string) error {
	url := s.base + `/write?db=` + s.db

	req, err := http.NewRequest("POST", url, strings.NewReader(line))
	if err != nil {
		return err
	}

	// add auth
	if s.username != "" && s.password != "" {
		req.SetBasicAuth(s.username, s.password)
	}

	_, err = s.c.Do(req)
	if err != nil {
		return err
	}

	return nil
}

// Ping will ping influxdb
func (s *SimpleInfluxDB) Ping() error {
	url := s.base + `/ping`

	req, err := http.NewRequest("GET", url, strings.NewReader(``))
	if err != nil {
		return err
	}

	// add auth
	if s.username != "" && s.password != "" {
		req.SetBasicAuth(s.username, s.password)
	}

	resp, err := s.c.Do(req)
	if err != nil {
		return err
	}

	log.Printf(`Ping to InfluxDB went HTTP %d`, resp.StatusCode)

	return nil
}
