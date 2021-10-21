package util

import (
	"gopkg.in/ini.v1"
	"log"
)

var (
	User     string
	Password string
	Url      string
)

func handleErr(err error) {
	if err != nil {
		log.Fatal("load config error: ", err)
	}
}

func LoadConfig() {
	cfg, err := ini.Load("config.ini")
	handleErr(err)
	sec := cfg.Section(ini.DefaultSection)
	user, err := sec.GetKey("user")
	handleErr(err)
	User = user.Value()

	pass, err := sec.GetKey("password")
	handleErr(err)
	Password = pass.Value()

	serve, err := cfg.GetSection("serve")
	handleErr(err)
	var (
		mqUser string
		mqPass string
		host   string
		port   string
		route  string
	)

	user, err = serve.GetKey("user")
	handleErr(err)
	mqUser = user.Value()

	pass, err = serve.GetKey("password")
	handleErr(err)
	mqPass = pass.Value()

	hostKey, err := serve.GetKey("host")
	handleErr(err)
	host = hostKey.Value()

	portKey, err := serve.GetKey("port")
	handleErr(err)
	port = portKey.Value()

	routeKey, err := serve.GetKey("route")
	handleErr(err)
	route = routeKey.Value()

	Url = "amqp://" + mqUser + ":" + mqPass + "@" + host + ":" + port + route
}
