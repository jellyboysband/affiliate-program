package main

import (
	"crawler/internal/core"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()

	log.Fatal(core.Run(log))
}
