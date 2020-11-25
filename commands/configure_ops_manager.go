package commands

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type ConfigureOpsManager struct {
	Config string `short:"c" help:"config file for the OpsManager" required`
}

func (ConfigureOpsManager) Run(cli *CLI) error {
	client := cli.newClient()

	contents, err := ioutil.ReadFile(cli.ConfigureOpsManager.Config)
	if err != nil {
		return err
	}

	var configs map[string]interface{}

	err = yaml.Unmarshal(contents, &configs)
	if err != nil {
		return err
	}

	return payloads.Update(client, configs)
}
