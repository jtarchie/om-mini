package commands

import (
	"crypto/tls"
	"fmt"
	"github.com/cloudfoundry-community/go-uaa"
	"github.com/go-resty/resty/v2"
	"golang.org/x/oauth2"
	"net/url"
)

type Credentials struct {
	Username             string   `env:"OM_USERNAME" short:"u" help:"admin username for the Ops Manager VM (not required for unauthenticated commands)"`
	Password             string   `env:"OM_PASSWORD" short:"p" help:"admin password for the Ops Manager VM (not required for unauthenticated commands)"`
	ClientID             string   `env:"OM_CLIENT_ID" help:"admin client ID for the Ops Manager VM (not required for unauthenticated commands)"`
	ClientSecret         string   `env:"OM_CLIENT_SECRET" help:"admin client secret for the Ops Manager VM (not required for unauthenticated commands)"`
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

func (c *CLI) newClient() *resty.Client {
	client := resty.New()

	if c.Target.Scheme == "" {
		c.Target.Scheme = "https"
	}

	client.SetHostURL(c.Target.String())

	var token *oauth2.Token

	client.SetRedirectPolicy(resty.FlexibleRedirectPolicy(10))
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: c.SkipSSLValidation})
	client.SetDebug(c.Verbose)

	client.OnBeforeRequest(func(_ *resty.Client, r *resty.Request) error {
		if token != nil && token.Valid() {
			r.Header.Set(
				"Authorization",
				fmt.Sprintf("Bearer %s", token.AccessToken),
			)
			return nil
		}

		options := []uaa.Option{
			uaa.WithSkipSSLValidation(c.SkipSSLValidation),
		}

		var authOption uaa.AuthenticationOption

		if c.Username != "" && c.Password != "" {
			authOption = uaa.WithPasswordCredentials(
				"opsman",
				"",
				c.Username,
				c.Password,
				uaa.OpaqueToken,
			)
		} else if c.ClientID != "" && c.ClientSecret != "" {
			authOption = uaa.WithClientCredentials(
				c.ClientID,
				c.ClientSecret,
				uaa.OpaqueToken,
			)
		}

		api, err := uaa.New(
			c.Target.String()+"/uaa",
			authOption,
			options...,
		)
		if err != nil {
			return fmt.Errorf("could not init UAA client: %w", err)
		}

		token, err = api.Token(r.Context())
		if err != nil {
			return fmt.Errorf("token could not be retrieved from target url: %w", err)
		}

		r.Header.Set(
			"Authorization",
			fmt.Sprintf("Bearer %s", token.AccessToken),
		)

		return nil
	})

	return client
}
