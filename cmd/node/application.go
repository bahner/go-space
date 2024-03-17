package main

import (
	"github.com/ergo-services/ergo/etf"
	"github.com/ergo-services/ergo/gen"
	log "github.com/sirupsen/logrus"
)

const (
	// Erlang application config
	VERSION = "0.0.1"
	NAME    = "go-space"
	DESC    = "SPACE node written in go to handle libp2p functionality."
)

type Application struct {
	gen.Application
}

func (app *Application) Load(args ...etf.Term) (gen.ApplicationSpec, error) {

	return gen.ApplicationSpec{
		Name:        NAME,
		Description: DESC,
		Version:     VERSION,
		Children: []gen.ApplicationChildSpec{
			{
				Name:  NAME,
				Child: new(SPACE),
			},
		},
	}, nil
}

func (app *Application) Start(process gen.Process, args ...etf.Term) {

	log.Infof("Application %s started with Pid %s", NAME, process.Self())
}
