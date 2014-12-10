package main

import (
	"encoding/json"
	"flag"
	"github.com/googollee/go-middleware/logger"
	"labix.org/v2/mgo"
	"net/http"
	"os"
)

type ErrorCallback interface {
	OnError(err error)
}

type Config struct {
	Mdb        string   `json:"mdb"`
	ListenAddr string   `json:"listen_addr"`
	App        string   `json:"app"`
	Ak         string   `json:"ak"`
	Sk         string   `json:"sk"`
	Assets     http.Dir `json:"assets"`
	Templates  string   `json:"templates"`

	refresh chan struct{}
	mdb     *mgo.Database
	logger  *logger.Logger
}

func Init() (*Config, error) {
	var confPath string
	flag.StringVar(&confPath, "conf", "./config.json", "configure file path")
	var listenAddr string
	flag.StringVar(&listenAddr, "listen", "", "listen address, overwrite config.json")
	var mdb string
	flag.StringVar(&mdb, "mdb", "", "mongo db address, overwrite config.json")

	flag.Parse()
	f, err := os.Open(confPath)
	if err != nil {
		return nil, err
	}
	var conf Config
	decoder := json.NewDecoder(f)
	if err := decoder.Decode(&conf); err != nil {
		return nil, err
	}

	if listenAddr != "" {
		conf.ListenAddr = listenAddr
	}
	if mdb != "" {
		conf.Mdb = mdb
	}

	conf.logger = logger.New("jipai", os.Stdout)
	conf.logger.Printf("connect mongo: %s", conf.Mdb)
	session, err := mgo.Dial(conf.Mdb)
	if err != nil {
		return nil, err
	}
	session.SetMode(mgo.Monotonic, true)
	conf.mdb = session.DB("")
	conf.refresh = make(chan struct{})

	go func() {
		for {
			<-conf.refresh
			session.Refresh()
		}
	}()

	return &conf, nil
}

func (c *Config) OnError(err error) {
	if err != nil {
		c.logger.Printf("mongo db error: %s", err)
	}
	select {
	case c.refresh <- struct{}{}:
	}
}
