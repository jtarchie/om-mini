package commands

import (
	"fmt"

	"github.com/jtarchie/om-mini/api"
	"gopkg.in/yaml.v3"
)

type StagedDirector struct{}

//nolint:gochecknoglobals
var directorPayloads = api.Payloads{
	"properties-configuration": api.Payload{
		Endpoint:      "/api/v0/staged/director/properties",
		IsCollectable: true,
	},
	"az-configuration": api.Payload{
		Endpoint:      "/api/v0/staged/director/availability_zones",
		Root:          "availability_zones",
		IsCollectable: true,
	},
}

func (StagedDirector) Run(cli *CLI) error {
	client := cli.newClient()

	configs, err := directorPayloads.Collect(client)
	if err != nil {
		return fmt.Errorf("could not collect director config: %w", err)
	}

	contents, err := yaml.Marshal(configs)
	if err != nil {
		return fmt.Errorf("could not marshal config file: %w", err)
	}

	fmt.Printf("%s", contents)

	return nil
}
