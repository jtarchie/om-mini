package main

import (
	"github.com/alecthomas/kong"
	"github.com/jtarchie/om-mini/commands"
)

func main() {
	cli := commands.CLI{}
	ctx := kong.Parse(&cli,
		kong.Name("om-mini"),
		kong.Description("om-mini helps you interact with an Ops Manager"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
		}),
	)
	err := ctx.Run(&cli)
	ctx.FatalIfErrorf(err)
}
