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

	appName := config.AppName
	description := config.Description
	version := config.Version

	return gen.ApplicationSpec{
		Name:        appName,
		Description: description,
		Version:     version,
		Children: []gen.ApplicationChildSpec{
			{
				Name:  appName,
				Child: createSpace(app.ctx),
			},
		},
	}, nil
}

func (app *Application) Start(process gen.Process, args ...etf.Term) {
	appName := config.AppName

	log.Infof("Application %s started with Pid %s\n", appName, process.Self())
}
