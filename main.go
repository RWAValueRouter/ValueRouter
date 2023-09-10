package main

import (
	"fmt"
	"github.com/anyswap/FastMulThreshold-DSA/cmd/utils"
	"github.com/anyswap/FastMulThreshold-DSA/log"
	"github.com/anyswap/FastMulThreshold-DSA/rpc"
	"github.com/anyswap/fastmpc-service-middleware/common"
	"github.com/anyswap/fastmpc-service-middleware/db"
	common2 "github.com/anyswap/fastmpc-service-middleware/internal/common"
	"github.com/anyswap/fastmpc-service-middleware/internal/flags"
	"github.com/anyswap/fastmpc-service-middleware/internal/params"
	_ "github.com/anyswap/fastmpc-service-middleware/jobs"
	service2 "github.com/anyswap/fastmpc-service-middleware/service"
	"gopkg.in/urfave/cli.v1"
	"net"
	"os"
	"os/signal"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"syscall"
)

const (
	clientIdentifier = "smw" // Client identifier to advertise over the network
)

var (
	// Git SHA1 commit hash of the release (set via linker flags)
	gitCommit  = ""
	gitDate    = ""
	gitVersion = ""
	// The app that holds all commands and flags.
	app      = flags.NewApp(gitCommit, gitDate, "the Smpc Wallet Service command line interface")
	rpcport  int
	endpoint string = "0.0.0.0"
	server   *rpc.Server

	logfile   string
	rotate    uint64
	maxage    uint64
	verbosity uint64
	json      bool
	color     bool

	stopLock   sync.Mutex
	signalChan = make(chan os.Signal, 1)
)

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func startRPCServer() error {
	go func() {
		server = rpc.NewServer()
		service := new(service2.ServiceMiddleWare)
		if err := server.RegisterName("smw", service); err != nil {
			log.Error("register service error", "error", err.Error())
			os.Exit(0)
		}

		// All APIs registered, start the HTTP listener
		var (
			listener net.Listener
			err      error
		)
		if rpcport == 0 {
			log.Error("rpc port is not specified")
			os.Exit(0)
			return
		}
		endpoint = endpoint + ":" + strconv.Itoa(rpcport)
		if listener, err = net.Listen("tcp", endpoint); err != nil {
			log.Error("listen tcp error", "error", err.Error())
			os.Exit(0)
			return
		}

		vhosts := common.SplitAndTrim("*")
		cors := common.SplitAndTrim("*")
		go func() {
			err2 := rpc.NewHTTPServer(cors, vhosts, rpc.DefaultHTTPTimeouts, server).Serve(listener)
			if err2 != nil {
				log.Error("============== new http server fail ==============", "err", err2)
				return
			}
		}()
		rpcstring := "==================== RPC Service Start! url = " + fmt.Sprintf("http://%s", endpoint) + " ====================="
		log.Info(rpcstring)

		exit := make(chan int)
		<-exit

		server.Stop()
	}()

	return nil
}

func StartSmw(c *cli.Context) {
	SetLogger()
	common.Init()
	db.Init()
	go func() {
		<-signalChan
		stopLock.Lock()
		common2.Info("=============================Cleaning before stop...======================================")
		stopLock.Unlock()
		os.Exit(0)
	}()
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := startRPCServer()
		if err != nil {
			log.Error("start rpc error" + err.Error())
			os.Exit(0)
		}
	}()

	select {}
}

// SetLogger config log print
func SetLogger() {
	common2.SetLogger(uint32(verbosity), json, color)
	if logfile != "" {
		common2.SetLogFile(logfile, rotate, maxage)
	}
}

func init() {
	app.Action = StartSmw
	app.HideVersion = true // we have a command to print the version
	app.Copyright = "Copyright 2018-2023 The anyswap Authors"
	app.Commands = []cli.Command{
		versionCommand,
		licenseCommand,
	}
	sort.Sort(cli.CommandsByName(app.Commands))
	app.Flags = []cli.Flag{
		cli.IntFlag{Name: "rpcport", Value: 0, Usage: "listen port", Destination: &rpcport},
		cli.StringFlag{Name: "logfile", Value: "", Usage: "Specify log file, support rotate", Destination: &logfile},
		cli.StringFlag{Name: "configfile", Value: "", Usage: "Specify config file, default config.json", Destination: &common.Configfile},
		cli.Uint64Flag{Name: "rotate", Value: 24, Usage: "log rotation time (unit hour)", Destination: &rotate},
		cli.Uint64Flag{Name: "maxage", Value: 7200, Usage: "log max age (unit hour)", Destination: &maxage},
		cli.Uint64Flag{Name: "verbosity", Value: 4, Usage: "log verbosity (0:panic, 1:fatal, 2:error, 3:warn, 4:info, 5:debug, 6:trace)", Destination: &verbosity},
		cli.BoolFlag{Name: "json", Usage: "output log in json format", Destination: &json},
		cli.BoolFlag{Name: "color", Usage: "output log in color text format", Destination: &color},
	}
	gitVersion = params.VersionWithMeta
}

var (
	versionCommand = cli.Command{
		Action:    utils.MigrateFlags(version),
		Name:      "version",
		Usage:     "Print version numbers",
		ArgsUsage: " ",
		Category:  "MISCELLANEOUS COMMANDS",
		Description: `
The output of this command is supposed to be machine-readable.
`,
	}
	licenseCommand = cli.Command{
		Action:    utils.MigrateFlags(license),
		Name:      "license",
		Usage:     "Display license information",
		ArgsUsage: " ",
		Category:  "MISCELLANEOUS COMMANDS",
	}
)

func version(ctx *cli.Context) error {
	fmt.Println(strings.Title(clientIdentifier))
	fmt.Println("Version:", params.VersionWithMeta)
	if gitCommit != "" {
		fmt.Println("Git Commit:", gitCommit)
	}
	if gitDate != "" {
		fmt.Println("Git Commit Date:", gitDate)
	}
	fmt.Println("Architecture:", runtime.GOARCH)
	fmt.Println("Protocol Versions:", params.ProtocolVersions)
	fmt.Println("Go Version:", runtime.Version())
	fmt.Println("Operating System:", runtime.GOOS)
	fmt.Printf("GOPATH=%s\n", os.Getenv("GOPATH"))
	fmt.Printf("GOROOT=%s\n", runtime.GOROOT())
	return nil
}

func license(_ *cli.Context) error {
	fmt.Println(`Copyright (C) 2018-2023  anyswap exchange Ltd. All rights reserved.
Copyright (C) 2018-2023 anyswap exchange

This library is free software; you can redistribute it and/or
modify it under the Apache License, Version 2.0.

This library is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.

See the License for the specific language governing permissions and
limitations under the License.`)
	return nil
}
