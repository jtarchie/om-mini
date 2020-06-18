package main

import (
	"crypto/tls"
	"fmt"
	"github.com/alecthomas/kong"
	"github.com/cloudfoundry-community/go-uaa"
	logger "github.com/izumin5210/gentleman-logger"
	"gopkg.in/h2non/gentleman.v2"
	"gopkg.in/h2non/gentleman.v2/context"
	gtls "gopkg.in/h2non/gentleman.v2/plugins/tls"
	"net/http"
	"os"
)

type Credentials struct {
	Username             string `env:"OM_USERNAME" short:"u" help:"admin username for the Ops Manager VM (not required for unauthenticated commands)"`
	Password             string `env:"OM_PASSWORD" short:"p" help:"admin password for the Ops Manager VM (not required for unauthenticated commands)"`
	Target               string `env:"OM_TARGET"   short:"t" required help:"location of the Ops Manager VM"`
	DecryptionPassphrase string `env:"OM_DECRYPTION_PASSPHRASE" short:"d" help:"Passphrase to decrypt the installation if the Ops Manager VM has been rebooted (optional for most commands)"`
	SkipSSLValidation    bool   `env:"OM_SKIP_SSL_VALIDATION" short:"k" help:"skip ssl certificate validation during http requests"`
}

type Curl struct {
	Data    string            `short:"d" help:"api request payload"`
	Header  map[string]string `short:"H" help:"used to specify custom headers with your command"`
	Path    string            `short:"p" help:"path to api endpoint" required`
	Request string            `short:"X" help:"http verb" default:"GET"`
	Verbose bool              `short:"v" help:"write all requests and responses to stderr"`
}

type CLI struct {
	Credentials

	Curl Curl `cmd`
}

func (Curl) Run(cli *CLI) error {
	client := gentleman.New()

	client.UseRequest(func(ctx *context.Context, h context.Handler) {
		transport, ok := ctx.Client.Transport.(*http.Transport)

		if cli.Username != "" && cli.Password != "" {
			api, err := uaa.New(
				cli.Target+"/uaa",
				uaa.WithPasswordCredentials("opsman", "", cli.Username, cli.Password, uaa.OpaqueToken),
				uaa.WithSkipSSLValidation(cli.SkipSSLValidation),
				uaa.WithVerbosity(cli.Curl.Verbose),
			)
			if err != nil {
				h.Error(ctx, err)
				return
			}

			token, err := api.Token(ctx)
			if err != nil {
				h.Error(ctx, err)
				return
			}

			ctx.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
		}
		h.Next(ctx)
	})

	client.Use(gtls.Config(&tls.Config{InsecureSkipVerify: cli.SkipSSLValidation}))

	if cli.Curl.Verbose {
		client.Use(logger.New(os.Stderr))
	}

	request := client.Request()
	request.
		Method(cli.Curl.Request).
		Path(cli.Curl.Path).
		BaseURL(cli.Target)

	if len(cli.Curl.Header) > 0 {
		request.SetHeaders(cli.Curl.Header)
	}

	response, err := request.Send()
	if err != nil {
		return err
	}

	if response.Ok {
		fmt.Println(response.String())
	}

	return nil
}



func main() {
	cli := CLI{}
	ctx := kong.Parse(&cli,
		kong.Name("om"),
		kong.Description("om helps you interact with an Ops Manager"),
		kong.UsageOnError(),
	)
	ctx.Bind()
	err := ctx.Run(&cli)
	ctx.FatalIfErrorf(err)
}
