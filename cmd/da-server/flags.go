package main

import (
	"encoding/hex"
	"errors"
	"fmt"

	opservice "github.com/ethereum-optimism/optimism/op-service"
	oplog "github.com/ethereum-optimism/optimism/op-service/log"
	sunrise "github.com/sunriselayer/sunrise-alt-da"
	"github.com/urfave/cli/v2"
)

const (
	ListenAddrFlagName       = "addr"
	PortFlagName             = "port"
	SunriseServerFlagName    = "sunrise.server"
	SunriseNamespaceFlagName = "sunrise.namespace"
)

const EnvVarPrefix = "OP_PLASMA_DA_SERVER"

func prefixEnvVars(name string) []string {
	return opservice.PrefixEnvVar(EnvVarPrefix, name)
}

var (
	ListenAddrFlag = &cli.StringFlag{
		Name:    ListenAddrFlagName,
		Usage:   "server listening address",
		Value:   "127.0.0.1",
		EnvVars: prefixEnvVars("ADDR"),
	}
	PortFlag = &cli.IntFlag{
		Name:    PortFlagName,
		Usage:   "server listening port",
		Value:   3100,
		EnvVars: prefixEnvVars("PORT"),
	}
	SunriseServerFlag = &cli.StringFlag{
		Name:    SunriseServerFlagName,
		Usage:   "sunrise server endpoint",
		Value:   "http://localhost:26658",
		EnvVars: prefixEnvVars("SUNRISE_SERVER"),
	}
	SunriseNamespaceFlag = &cli.StringFlag{
		Name:    SunriseNamespaceFlagName,
		Usage:   "sunrise namespace",
		Value:   "",
		EnvVars: prefixEnvVars("SUNRISE_NAMESPACE"),
	}
)

var requiredFlags = []cli.Flag{
	ListenAddrFlag,
	PortFlag,
}

var optionalFlags = []cli.Flag{
	SunriseServerFlag,
	SunriseNamespaceFlag,
}

func init() {
	optionalFlags = append(optionalFlags, oplog.CLIFlags(EnvVarPrefix)...)
	Flags = append(requiredFlags, optionalFlags...)
}

// Flags contains the list of configuration options available to the binary.
var Flags []cli.Flag

type CLIConfig struct {
	SunriseEndpoint  string
	SunriseNamespace string
}

func ReadCLIConfig(ctx *cli.Context) CLIConfig {
	return CLIConfig{
		SunriseEndpoint:  ctx.String(SunriseServerFlagName),
		SunriseNamespace: ctx.String(SunriseNamespaceFlagName),
	}
}

func (c CLIConfig) Check() error {
	if c.SunriseEnabled() && (c.SunriseEndpoint == "" || c.SunriseNamespace == "") {
		return errors.New("all Sunrise flags must be set")
	}
	if c.SunriseEnabled() {
		if _, err := hex.DecodeString(c.SunriseNamespace); err != nil {
			return err
		}
	}
	return nil
}

func (c CLIConfig) SunriseConfig() sunrise.SunriseConfig {
	ns, _ := hex.DecodeString(c.SunriseNamespace)
	return sunrise.SunriseConfig{
		URL:       c.SunriseEndpoint,
		Namespace: ns,
	}
}

func (c CLIConfig) SunriseEnabled() bool {
	return !(c.SunriseEndpoint == "" && c.SunriseNamespace == "")
}

func CheckRequired(ctx *cli.Context) error {
	for _, f := range requiredFlags {
		if !ctx.IsSet(f.Names()[0]) {
			return fmt.Errorf("flag %s is required", f.Names()[0])
		}
	}
	return nil
}
