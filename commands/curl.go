package commands

import "fmt"

type Curl struct {
	Data    string            `short:"d" help:"api request payload"`
	Header  map[string]string `short:"H" help:"used to specify custom headers with your command"`
	Path    string            `short:"p" help:"path to api endpoint" required`
	Request string            `short:"X" help:"http verb" default:"GET"`
}

func (Curl) Run(cli *CLI) error {
	client := cli.newClient()

	request := client.Request()
	request.
		Method(cli.Curl.Request).
		Path(cli.Curl.Path)

	if len(cli.Curl.Header) > 0 {
		request.SetHeaders(cli.Curl.Header)
	}

	if cli.Curl.Data != "" {
		request.BodyString(cli.Curl.Data)
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
