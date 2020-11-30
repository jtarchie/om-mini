package api

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
	"os"
)

type Payload struct {
	Endpoint      string
	Root          string
	IsCollectable bool
}

type Payloads map[string]Payload
type configs map[string]interface{}

func (p Payloads) Update(
	client *resty.Client,
	configs configs,
) error {
	for setting, payload := range p {
		if payload.Endpoint == "" {
			continue
		}

		config, ok := configs[setting]
		if !ok {
			continue
		}

		_, _ = fmt.Fprintf(os.Stderr, "updating %s\n", setting)

		request := client.R()

		if payload.Root == "" {
			request.SetBody(config)
		} else {
			request.SetBody(map[string]interface{}{
				payload.Root: config,
			})
		}

		response, err := request.Execute(
			http.MethodPut,
			payload.Endpoint,
		)
		if err != nil {
			return err
		}

		if !response.IsSuccess() {
			return fmt.Errorf(
				"could not update %q, the Ops Manager API returned an error:\n%s",
				setting,
				response.String(),
			)
		}
	}

	return nil
}

func (p Payloads) Collect(client *resty.Client) (configs, error) {
	configs := configs{}

	for setting, payload := range p {
		if !payload.IsCollectable {
			_, _ = fmt.Fprintf(os.Stderr, "unable to collect %q, skipping\n", setting)
			continue
		}

		if payload.Endpoint == "" {
			continue
		}

		_, _ = fmt.Fprintf(os.Stderr, "collecting %q\n", setting)

		request := client.R()

		response, err := request.Execute(
			http.MethodGet,
			payload.Endpoint,
		)
		if err != nil {
			return configs, err
		}

		if !response.IsSuccess() {
			return configs, fmt.Errorf(
				"could not collect %q, the Ops Manager API returned an error:\n%s",
				setting,
				response.String(),
			)
		}

		config := map[string]interface{}{}

		err = json.Unmarshal(response.Body(), &config)
		if err != nil {
			return configs, fmt.Errorf(
				"could unmarshal %q from Ops Manager API: %w",
				setting,
				err,
			)
		}

		if payload.Root == "" {
			configs[setting] = config
		} else {
			configs[setting] = config[payload.Root]
		}

	}

	return configs, nil
}
