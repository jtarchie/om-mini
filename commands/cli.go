package commands

import (
	"crypto/tls"
	"fmt"
	"github.com/cloudfoundry-community/go-uaa"
	logger "github.com/izumin5210/gentleman-logger"
	"golang.org/x/oauth2"
	"gopkg.in/h2non/gentleman.v2"
	"gopkg.in/h2non/gentleman.v2/context"
	"gopkg.in/h2non/gentleman.v2/plugins/redirect"
	gtls "gopkg.in/h2non/gentleman.v2/plugins/tls"
	"net/url"
	"os"
)

type Credentials struct {
	Username             string   `env:"OM_USERNAME" short:"u" help:"admin username for the Ops Manager VM (not required for unauthenticated commands)"`
	Password             string   `env:"OM_PASSWORD" short:"p" help:"admin password for the Ops Manager VM (not required for unauthenticated commands)"`
	Target               *url.URL `env:"OM_TARGET"   short:"t" required help:"location of the Ops Manager VM"`
	DecryptionPassphrase string   `env:"OM_DECRYPTION_PASSPHRASE" short:"d" help:"Passphrase to decrypt the installation if the Ops Manager VM has been rebooted (optional for most commands)"`
	SkipSSLValidation    bool     `env:"OM_SKIP_SSL_VALIDATION" short:"k" help:"skip ssl certificate validation during http requests"`
	Verbose              bool     `short:"v" help:"write all requests and responses to stderr"`
}

type CLI struct {
	Credentials

	Curl                Curl                `cmd`
	ConfigureOpsManager ConfigureOpsManager `cmd`
	StagedOpsManager    StagedOpsManager    `cmd`
}

func (c *CLI) newClient() *gentleman.Client {
	client := gentleman.New()

	if c.Target.Scheme == "" {
		c.Target.Scheme = "https"
	}

	client.BaseURL(c.Target.String())

	var token *oauth2.Token

	client.Use(redirect.Config(redirect.Options{
		Limit:   10,
		Trusted: true,
	}))

	client.UseRequest(func(ctx *context.Context, h context.Handler) {
		if token != nil && token.Valid() {
			ctx.Request.Header.Set(
				"Authorization",
				fmt.Sprintf("Bearer %s", token.AccessToken),
			)
			h.Next(ctx)
			return
		}

		if c.Username != "" && c.Password != "" {
			api, err := uaa.New(
				c.Target.String()+"/uaa",
				uaa.WithPasswordCredentials("opsman", "", c.Username, c.Password, uaa.OpaqueToken),
				uaa.WithSkipSSLValidation(c.SkipSSLValidation),
				uaa.WithVerbosity(c.Verbose),
			)
			if err != nil {
				h.Error(ctx, err)
				return
			}

			token, err = api.Token(ctx)
			if err != nil {
				h.Error(ctx, err)
				return
			}

			ctx.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
		}
		h.Next(ctx)
	})

	client.Use(gtls.Config(&tls.Config{InsecureSkipVerify: c.SkipSSLValidation}))

	if c.Verbose {
		client.Use(logger.New(os.Stderr))
	}

	return client
}
