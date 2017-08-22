package main

// Geeny CLI
// @see https://developers.geeny.io/documentation/cli/

import (
	"fmt"
	"os"

	"geeny/cli"
	"geeny/cli/tree"
	"geeny/config"
	"geeny/dyn"
	log "geeny/log"
	"geeny/output"
	"geeny/plugin"
	"geeny/version"
)

func main() {
	commandLine := cli.NewCommandLine()
	isBashCompletion := commandLine.Search(os.Args, cli.FlagGenerateBashCompletion)
	if isBashCompletion {
		log.Set(false)
		commandLine.SetActionsEnabled(false) // disable actions if found a bash completion flag
		config.CurrentInt.SpinnerOutput = false
		config.CurrentExt.Log = false
		config.CurrentExt.SwaggerValidate = false
		config.CurrentExt.AutoUpdateCheck = false
	}
	log.Info("Main started")
	log.Tracef("User configuration: %+v", config.CurrentExt)

	// check for new version after running the cli, in case bash completion is passed
	if config.CurrentExt.AutoUpdateCheck && !isBashCompletion {
		newVersion, message := version.CheckUpdate(true)
		if newVersion {
			fmt.Fprintln(os.Stderr, message)
		}
	}

	// support undocumented api endpoint override for testing purposes
	var apiUrl = os.Getenv("GEENY_API_URL")
	if len(apiUrl) == 0 {
		apiUrl = config.CurrentExt.ApiUrl
	}
	var connectUrl = os.Getenv("GEENY_CONNECT_URL")
	if len(connectUrl) == 0 {
		connectUrl = config.CurrentExt.ConnectUrl
	}

	// initialize command tree
	tree := tree.NewCommandTree(apiUrl, connectUrl)
	tree.SetupParents()
	//TODO only initialize plugins until the command to run is in the tree
	log.Info("Initializing plugins")
	var p plugin.Plugin
	var err error
	p, err = initPlugin(connectUrl, config.CurrentExt.ConnectSwaggerUrl, tree, isBashCompletion)
	tree.SetupParents()
	defer p.Close()
	p, err = initPlugin(apiUrl, config.CurrentExt.ApiSwaggerUrl, tree, isBashCompletion)
	defer p.Close()
	tree.SetupParents()
	tree.Sort()

	// run cli
	_, err = commandLine.Run(tree, os.Args)
	checkError(err, isBashCompletion)

	log.Info("Main finished")
}

func initPlugin(url string, swaggerUrl string, tree *cli.Command, isBashCompletion bool) (p plugin.Plugin, err error) {
	p = dyn.NewPlugin(url, swaggerUrl) //TODO make this a real plugin
	defer p.Close()
	err = p.Init(tree)
	checkError(err, isBashCompletion)
	return p, err
}

func checkError(err error, isBashCompletion bool) {
	if err != nil {
		if len(err.Error()) > 0 && !isBashCompletion {
			output.Println(err)
			log.Fatalf("Main error: %v", err)
		} else {
			os.Exit(1)
		}
	}
}
