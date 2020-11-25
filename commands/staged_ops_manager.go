package commands

import (
	"fmt"
	"github.com/jtarchie/om-mini/api"
	"gopkg.in/yaml.v3"
)

type StagedOpsManager struct{}

var payloads = api.Payloads{
	"syslog-settings": api.Payload{
		Endpoint:      "/api/v0/settings/syslog",
		Root:          "syslog",
		IsCollectable: true,
	},
	"banner-settings": api.Payload{
		Endpoint:      "/api/v0/settings/banner",
		IsCollectable: true,
	},
	"rbac-settings": api.Payload{
		Endpoint:      "/api/v0/settings/rbac",
		IsCollectable: false,
	},
	"pivotal-network-settings": api.Payload{
		Endpoint:      "/api/v0/settings/pivotal_network_settings",
		Root:          "pivotal_network_settings",
		IsCollectable: false,
	},
	"ssl-certificate": api.Payload{
		Endpoint:      "/api/v0/settings/ssl_certificate",
		IsCollectable: true,
	},
	"uaa-tokens-expiration": api.Payload{
		Endpoint:      "/api/v0/uaa/tokens_expiration",
		IsCollectable: true,
		Root:          "tokens_expiration",
	},
}

func (StagedOpsManager) Run(cli *CLI) error {
	client := cli.newClient()

	configs, err := payloads.Collect(client)
	if err != nil {
		return err
	}

	contents, err := yaml.Marshal(configs)
	if err != nil {
		return err
	}

	fmt.Printf("%s", contents)

	return nil
}
