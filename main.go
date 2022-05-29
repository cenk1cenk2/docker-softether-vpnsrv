package main

import (
	"github.com/urfave/cli/v2"

	pipe "github.com/cenk1cenk2/docker-softether-vpnsrv/pipe"
	. "gitlab.kilic.dev/libraries/plumber/v2"
)

func main() {
	p := Plumber{}

	p.New(
		func(p *Plumber) *cli.App {
			return &cli.App{
				Name:        CLI_NAME,
				Version:     VERSION,
				Usage:       DESCRIPTION,
				Description: DESCRIPTION,
				Flags:       p.AppendFlags(pipe.Flags),
				Action: func(ctx *cli.Context) error {
					return pipe.TL.RunJobs(
						pipe.TL.JobSequence(
							pipe.New(p).Job(ctx),
						),
					)
				},
			}
		}).Run()
}
