// +build headless

package app

import (
	"os"

	"github.com/urfave/cli"

	"github.com/p9c/pod/pkg/conte"
)

var guiHandle = func(cx *conte.Xt) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		Warn("GUI was disabled for this build (server only version)")
		os.Exit(1)
		return nil
	}
}

var monitorHandle = func(cx *conte.Xt) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		Warn("GUI was disabled for this build (server only version)")
		os.Exit(1)
		return nil
	}
}
