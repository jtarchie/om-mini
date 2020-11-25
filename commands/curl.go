package commands

import "fmt"

type Curl struct {
	Data   string            `short:"d" help:"api request payload"`
	Header map[string]string `short:"H" help:"used to specify custom headers with your command"`
	Path   string            `short:"p" help:"path to api endpoint" required`
	Method string            `short:"X" help:"http verb" default:"GET"`
}

func (Curl) Run(cli *CLI) error {
	client := cli.newClient()

	request := client.R()

	if len(cli.Curl.Header) > 0 {
		request.SetHeaders(cli.Curl.Header)
	}

	if cli.Curl.Data != "" {
		request.SetBody(cli.Curl.Data)
	}

	response, err := request.Execute(
		cli.Curl.Method,
		cli.Curl.Path,
	)
	if err != nil {
		return err
	}

	if response.IsSuccess() {
		fmt.Println(response.String())
	}

	return nil
}
