package app

import (
	"github.com/bahner/go-space/config"
	"github.com/ergo-services/ergo/etf"
	"github.com/ergo-services/ergo/gen"
	log "github.com/sirupsen/logrus"
)

type Application struct {
	gen.Application
}

func (app *Application) Load(args ...etf.Term) (gen.ApplicationSpec, error) {

	return gen.ApplicationSpec{
		Name:        config.NAME,
		Description: config.DESC,
		Version:     config.VERSION,
		Children: []gen.ApplicationChildSpec{
			{
				Name:  config.NAME,
				Child: new(SPACE),
			},
		},
	}, nil
}

func (app *Application) Start(process gen.Process, args ...etf.Term) {
	appName := config.NAME

	log.Infof("Application %s started with Pid %s", appName, process.Self())
}
