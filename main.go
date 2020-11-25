package main

import (
	"github.com/alecthomas/kong"
	"github.com/jtarchie/om-mini/commands"
)

func main() {
	cli := commands.CLI{}
	ctx := kong.Parse(&cli,
		kong.Name("om"),
		kong.Description("om helps you interact with an Ops Manager"),
		kong.UsageOnError(),
	)
	ctx.Bind()
	err := ctx.Run(&cli)
	ctx.FatalIfErrorf(err)
}
