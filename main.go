package main

import (
	"def2sql/config"
	"def2sql/service"

	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()

	conf, err := config.Read("config.yml")
	if err != nil {
		log.Fatal(err)
	}

	s, err := service.New(conf)
	if err != nil {
		log.Fatal(err)
	}

	err = s.Run()
	if err != nil {
		log.Errorf("%+v", err)
	}
}
