package main

import (
	"github.com/urfave/cli/v2"

	pipe "github.com/cenk1cenk2/docker-softether-vpnsrv/pipe"
	. "gitlab.kilic.dev/libraries/plumber/v3"
)

func main() {
	p := Plumber{
		DocsFile:         "CLI.md",
		DocsExcludeFlags: true,
	}

	p.New(
		func(p *Plumber) *cli.App {
			return &cli.App{
				Name:        CLI_NAME,
				Version:     VERSION,
				Usage:       DESCRIPTION,
				Description: DESCRIPTION,
				Flags:       p.AppendFlags(pipe.Flags),
				Before: func(ctx *cli.Context) error {
					p.EnableTerminator()

					return nil
				},
				Action: func(ctx *cli.Context) error {
					return pipe.TL.RunJobs(
						pipe.New(p).SetCliContext(ctx).Job(),
					)
				},
			}
		}).
		Run()
}
