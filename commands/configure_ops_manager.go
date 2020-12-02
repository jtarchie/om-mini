package commands

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type ConfigureOpsManager struct {
	Config string `short:"c" help:"config file for the OpsManager" required:""`
}

func (ConfigureOpsManager) Run(cli *CLI) error {
	client := cli.newClient()

	contents, err := ioutil.ReadFile(cli.ConfigureOpsManager.Config)
	if err != nil {
		return fmt.Errorf("could not read config file: %w", err)
	}

	var configs map[string]interface{}

	err = yaml.Unmarshal(contents, &configs)
	if err != nil {
		return fmt.Errorf("could not unmarshal config file: %w", err)
	}

	return opsmanagerPayloads.Update(client, configs)
}
