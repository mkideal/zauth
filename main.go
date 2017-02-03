package main

import (
	"os"

	"github.com/mkideal/cli"
	clix "github.com/mkideal/cli/ext"
	"github.com/mkideal/log"
	"github.com/mkideal/log/logger"
	"github.com/mkideal/pkg/osutil/signal"

	"bitbucket.org/mkideal/accountd/etc"
	"bitbucket.org/mkideal/accountd/model"
	"bitbucket.org/mkideal/accountd/server"
)

type argT struct {
	cli.Helper
	Version      bool         `cli:"!v,version" usage:"display version"`
	PidFile      clix.PidFile `cli:"pid" usage:"pid filepath" dft:"./var/accountd.pid"`
	LogLevel     logger.Level `cli:"log-level" usage:"log level: trace/debug/info/warn/error/fatal" dft:"info"`
	LogProviders string       `cli:"log-providers" usage:"log providers seperated by /" dft:"colored_console/file"`
	LogOpts      string       `cli:"log-opts" usage:"log options formatted with json or form" dft:"dir=./var/logs"`
	server.Config
}

const successPrefix = "accountd start ok"

var root = &cli.Command{
	Name: "accountd",
	Desc: "account service",
	Argv: func() interface{} { return new(argT) },

	Fn: func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)
		if argv.Version {
			printVersion(ctx)
			return nil
		}

		// initialize log
		log.Init(argv.LogProviders, argv.LogOpts)
		log.SetLevel(argv.LogLevel)

		log.WithJSON(argv).Debug("argv")

		// check pid file
		if err := argv.PidFile.New(); err != nil {
			log.Error("Error: %v", err)
			return err
		}
		defer argv.PidFile.Remove()

		// initialize model
		if err := model.Init(); err != nil {
			log.Error("Error: %v", err)
			return err
		}

		// run server
		svr := server.New(argv.Config)
		if err := svr.Run(); err != nil {
			return err
		}
		// quit server
		defer svr.Quit()

		// notify daemon's parent process
		cli.DaemonResponse(successPrefix)

		// block until signal INT received
		signal.Wait(os.Interrupt)

		return nil
	},
}

var daemon = &cli.Command{
	Name: "daemon",
	Desc: "startup service as a daemon process",
	Argv: func() interface{} { return new(argT) },

	Fn: func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)
		if argv.Version {
			printVersion(ctx)
			return nil
		}
		return cli.Daemon(ctx, successPrefix)
	},
}

func main() {
	defer log.Uninit(nil)
	cli.Root(root,
		cli.Tree(daemon),
	).Run(os.Args[1:])
}

func printVersion(ctx *cli.Context) {
	ctx.String("v" + etc.Version + "\n")
}
