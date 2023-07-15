package commands

import (
	"github.com/simonwep/genisis/core"
	"github.com/simonwep/genisis/routes"
	"github.com/urfave/cli/v2"
)

func Start(*cli.Context) error {
	router := routes.SetupRoutes()

	if err := router.SetTrustedProxies(nil); err != nil {
		return err
	} else if err = router.Run("0.0.0.0:" + core.Config.AppPort); err != nil {
		return err
	}

	return nil
}
