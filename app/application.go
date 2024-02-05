package app

import (
	"context"

	"github.com/bahner/go-space/config"
	"github.com/ergo-services/ergo/etf"
	"github.com/ergo-services/ergo/gen"
	log "github.com/sirupsen/logrus"
)

func createApplication(ctx context.Context) gen.ApplicationBehavior {
	return &Application{
		ctx: ctx,
	}
}

type Application struct {
	gen.Application
	ctx context.Context
}

func (app *Application) Load(args ...etf.Term) (gen.ApplicationSpec, error) {

	return gen.ApplicationSpec{
		Name:        config.NAME,
		Description: config.DESC,
		Version:     config.VERSION,
		Children: []gen.ApplicationChildSpec{
			{
				Name:  config.NAME,
				Child: createSPACE(app.ctx),
			},
		},
	}, nil
}

func (app *Application) Start(process gen.Process, args ...etf.Term) {
	appName := config.NAME

	log.Infof("Application %s started with Pid %s", appName, process.Self())
}
