package main

import (
	"encoding/hex"
	"errors"
	"fmt"

	opservice "github.com/ethereum-optimism/optimism/op-service"
	oplog "github.com/ethereum-optimism/optimism/op-service/log"
	"github.com/urfave/cli/v2"

	sunrise "github.com/sunriselayer/sunrise-alt-da"
)

const (
	ListenAddrFlagName    = "addr"
	PortFlagName          = "port"
	GenericCommFlagName   = "generic-commitment"
	SunriseServerFlagName = "sunrise.server"
	// SunriseAuthTokenFlagName = "sunrise.auth-token"
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
	GenericCommFlag = &cli.BoolFlag{
		Name:    GenericCommFlagName,
		Usage:   "enable generic commitments for testing. Not for production use.",
		EnvVars: prefixEnvVars("GENERIC_COMMITMENT"),
		Value:   true,
	}
	SunriseServerFlag = &cli.StringFlag{
		Name:    SunriseServerFlagName,
		Usage:   "sunrise server endpoint",
		Value:   "http://localhost:26658",
		EnvVars: prefixEnvVars("SUNRISE_SERVER"),
	}
	// SunriseAuthTokenFlag = &cli.StringFlag{
	// 	Name:    SunriseAuthTokenFlagName,
	// 	Usage:   "sunrise auth token",
	// 	Value:   "",
	// 	EnvVars: prefixEnvVars("SUNRISE_AUTH_TOKEN"),
	// }
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
	GenericCommFlag,
	SunriseServerFlag,
	// SunriseAuthTokenFlag,
	SunriseNamespaceFlag,
}

func init() {
	optionalFlags = append(optionalFlags, oplog.CLIFlags(EnvVarPrefix)...)
	Flags = append(requiredFlags, optionalFlags...)
}

// Flags contains the list of configuration options available to the binary.
var Flags []cli.Flag

type CLIConfig struct {
	UseGenericComm  bool
	SunriseEndpoint string
	// SunriseAuthToken string
	SunriseNamespace string
}

func ReadCLIConfig(ctx *cli.Context) CLIConfig {
	return CLIConfig{
		UseGenericComm:  ctx.Bool(GenericCommFlagName),
		SunriseEndpoint: ctx.String(SunriseServerFlagName),
		// SunriseAuthToken: ctx.String(SunriseAuthTokenFlagName),
		SunriseNamespace: ctx.String(SunriseNamespaceFlagName),
	}
}

func (c CLIConfig) Check() error {
	if c.SunriseEnabled() && (c.SunriseEndpoint == "" || /* c.SunriseAuthToken == "" || */ c.SunriseNamespace == "") {
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
		URL: c.SunriseEndpoint,
		// AuthToken: c.SunriseAuthToken,
		Namespace: ns,
	}
}

func (c CLIConfig) SunriseEnabled() bool {
	return !(c.SunriseEndpoint == "" && /* c.SunriseAuthToken == "" && */ c.SunriseNamespace == "")
}

func CheckRequired(ctx *cli.Context) error {
	for _, f := range requiredFlags {
		if !ctx.IsSet(f.Names()[0]) {
			return fmt.Errorf("flag %s is required", f.Names()[0])
		}
	}
	return nil
}
