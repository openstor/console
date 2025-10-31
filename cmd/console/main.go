// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/minio/cli"
	"github.com/openstor/console/pkg"
	"github.com/openstor/pkg/v3/console"
	"github.com/openstor/pkg/v3/trie"
	"github.com/openstor/pkg/v3/words"
)

// Help template for Console.
var consoleHelpTemplate = `NAME:
 {{.Name}} - {{.Usage}}

DESCRIPTION:
 {{.Description}}

USAGE:
 {{.HelpName}} {{if .VisibleFlags}}[FLAGS] {{end}}COMMAND{{if .VisibleFlags}}{{end}} [ARGS...]

COMMANDS:
 {{range .VisibleCommands}}{{join .Names ", "}}{{ "\t" }}{{.Usage}}
 {{end}}{{if .VisibleFlags}}
FLAGS:
 {{range .VisibleFlags}}{{.}}
 {{end}}{{end}}
VERSION:
 {{.Version}}
`

func newApp(name string) *cli.App {
	// Collection of console commands currently supported are.
	var commands []cli.Command

	// Collection of console commands currently supported in a trie tree.
	commandsTree := trie.NewTrie()

	// registerCommand registers a cli command.
	registerCommand := func(command cli.Command) {
		commands = append(commands, command)
		commandsTree.Insert(command.Name)
	}

	// register commands
	for _, cmd := range appCmds {
		registerCommand(cmd)
	}

	findClosestCommands := func(command string) []string {
		var closestCommands []string
		closestCommands = append(closestCommands, commandsTree.PrefixMatch(command)...)

		sort.Strings(closestCommands)
		// Suggest other close commands - allow missed, wrongly added and
		// even transposed characters
		for _, value := range commandsTree.Walk(commandsTree.Root()) {
			if sort.SearchStrings(closestCommands, value) < len(closestCommands) {
				continue
			}
			// 2 is arbitrary and represents the max
			// allowed number of typed errors
			if words.DamerauLevenshteinDistance(command, value) < 2 {
				closestCommands = append(closestCommands, value)
			}
		}

		return closestCommands
	}

	cli.HelpFlag = cli.BoolFlag{
		Name:  "help, h",
		Usage: "show help",
	}

	app := cli.NewApp()
	app.Name = name
	app.Version = pkg.Version + " - " + pkg.ShortCommitID
	app.Author = "MinIO, Inc."
	app.Usage = "OpenStor Console Server"
	app.Description = `OpenStor Console Server`
	app.Copyright = "(c) 2021 OpenStor contributors"
	app.Compiled, _ = time.Parse(time.RFC3339, pkg.ReleaseTime)
	app.Commands = commands
	app.HideHelpCommand = true // Hide `help, h` command, we already have `minio --help`.
	app.CustomAppHelpTemplate = consoleHelpTemplate
	app.CommandNotFound = func(_ *cli.Context, command string) {
		console.Printf("‘%s’ is not a console sub-command. See ‘console --help’.\n", command)
		closestCommands := findClosestCommands(command)
		if len(closestCommands) > 0 {
			console.Println()
			console.Println("Did you mean one of these?")
			for _, cmd := range closestCommands {
				console.Printf("\t‘%s’\n", cmd)
			}
		}

		os.Exit(1)
	}

	return app
}

func main() {
	args := os.Args
	// Set the orchestrator app name.
	appName := filepath.Base(args[0])
	// Run the app - exit on error.
	if err := newApp(appName).Run(args); err != nil {
		os.Exit(1)
	}
}
