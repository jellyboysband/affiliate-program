package core

import (
	"crawler/internal/modules/aliexpress"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type module func(log *logrus.Logger) error

func Run(log *logrus.Logger) error {

	OnModule(log, aliexpress.Name, aliexpress.RunCollector)

	return nil
}

func OnModule(log *logrus.Logger, moduleName string, module module) {
	for {
		err := module(log)
		log.Warn(errors.Wrapf(err, "failed to module %s", moduleName))
	}
}
