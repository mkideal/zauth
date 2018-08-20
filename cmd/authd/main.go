package main

import (
	"os"

	"github.com/mkideal/cli"
	clix "github.com/mkideal/cli/ext"
	"github.com/mkideal/log"
	"github.com/mkideal/log/logger"
	"github.com/mkideal/pkg/osutil/signal"
	"github.com/mkideal/pkg/service/discovery"

	"github.com/mkideal/accountd/etc"
	"github.com/mkideal/accountd/server"
)

type argT struct {
	cli.Helper
	Version       bool         `cli:"!v,version" usage:"display version" json:"-"`
	PidFile       clix.PidFile `cli:"pid" usage:"pid filepath" dft:"./var/authd.pid"`
	LogLevel      logger.Level `cli:"log-level" usage:"log level: trace/debug/info/warn/error/fatal" dft:"info"`
	LogProviders  string       `cli:"log-providers" usage:"log providers seperated by /" dft:"colored_console/file"`
	LogOpts       string       `cli:"log-opts" usage:"log options formatted with json or form" dft:"dir=./var/logs"`
	EtcdEndpoints string       `cli:"etcd" usage:"etcd endpoints for service register"`
	ServiceName   string       `cli:"service-name" usage:"registered service name" dft:"authd"`
	server.Config
}

const successPrefix = "authd start ok"

var root = &cli.Command{
	Name: "authd",
	Desc: "authorization service",
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
			cli.DaemonResponse(err.Error())
			return err
		}
		defer argv.PidFile.Remove()

		// run server
		svr, err := server.New(argv.Config)
		if err != nil {
			log.Error("Error: %v", err)
			cli.DaemonResponse(err.Error())
			return err
		}
		if err := svr.Run(); err != nil {
			log.Error("Error: %v", err)
			cli.DaemonResponse(err.Error())
			return err
		}

		// register service
		var discoveryClient *discovery.Discovery
		if argv.EtcdEndpoints != "" {
			discoveryClient = &discovery.Discovery{EtcdEndpoints: argv.EtcdEndpoints}
			if err := discoveryClient.Init(); err != nil {
				log.Error("Error: init discovery: %v", err)
				cli.DaemonResponse(err.Error())
				return err
			}
			discovery.Interval(*discoveryClient, func(ttl *discovery.TTL) {
				address := discovery.Address{Addr: argv.Addr}
				discoveryClient.Register(argv.ServiceName, address, ttl.Opt())
			}, discovery.DefaultTTL)
		}
		// quit server
		defer func() {
			if discoveryClient != nil {
				// unregister service
				discoveryClient.Unregister(argv.ServiceName, argv.Addr)
			}
			svr.Quit()
		}()

		// notify daemon's parent process
		cli.DaemonResponse(successPrefix)

		// block until signal INT received
		signal.Wait(os.Interrupt)

		return nil
	},
}

var daemon = &cli.Command{
	Name: "daemon",
	Desc: "startup service as a background process",
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
	defer log.Uninit(log.InitConsole(log.LvWARN))
	err := cli.Root(root,
		cli.Tree(daemon),
	).Run(os.Args[1:])
	if err != nil {
		log.Error("Error: %v", err)
	}
}

func printVersion(ctx *cli.Context) {
	ctx.String("v" + etc.Version + "\n")
}
