package app

import (
	"context"

	"github.com/ergo-services/ergo/etf"
	"github.com/ergo-services/ergo/gen"
)

func StartApplication(ctx context.Context) gen.ApplicationBehavior {
	return &Application{
		ctx: ctx,
	}
}

type Application struct {
	gen.Application
	ctx context.Context
}

func (app *Application) Load(args ...etf.Term) (gen.ApplicationSpec, error) {

	nodeInit(app.ctx)

	return gen.ApplicationSpec{
		Name:        appName,
		Description: description,
		Version:     version,
		Children: []gen.ApplicationChildSpec{
			{
				Name:  appName,
				Child: createMyspace(app.ctx),
			},
		},
	}, nil
}

func (app *Application) Start(process gen.Process, args ...etf.Term) {
	log.Infof("Application %s started with Pid %s\n", appName, process.Self())
}
