package main

import (
	"fmt"

	oplog "github.com/ethereum-optimism/optimism/op-service/log"
	"github.com/ethereum-optimism/optimism/op-service/opio"
	"github.com/urfave/cli/v2"

	sunrise "github.com/sunriselayer/sunrise-alt-da"
)

type Server interface {
	Start() error
	Stop() error
}

func StartDAServer(cliCtx *cli.Context) error {
	if err := CheckRequired(cliCtx); err != nil {
		return err
	}

	cfg := ReadCLIConfig(cliCtx)
	if err := cfg.Check(); err != nil {
		return err
	}

	logCfg := oplog.ReadCLIConfig(cliCtx)

	l := oplog.NewLogger(oplog.AppOut(cliCtx), logCfg)
	oplog.SetGlobalLogHandler(l.Handler())

	l.Info("Initializing Plasma DA server...")

	var server Server

	switch {
	case cfg.SunriseEnabled():
		l.Info("Using sunrise storage", "url", cfg.SunriseConfig().URL)
		store := sunrise.NewSunriseStore(cfg.SunriseConfig())
		server = sunrise.NewSunriseServer(cliCtx.String(ListenAddrFlagName), cliCtx.Int(PortFlagName), store, l)
	}

	if err := server.Start(); err != nil {
		return fmt.Errorf("failed to start the DA server")
	} else {
		l.Info("Started DA Server")
	}

	defer func() {
		if err := server.Stop(); err != nil {
			l.Error("failed to stop DA server", "err", err)
		}
	}()

	opio.BlockOnInterrupts()

	return nil
}
