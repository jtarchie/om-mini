package api

import (
	"fmt"
	"gopkg.in/h2non/gentleman.v2"
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
	client *gentleman.Client,
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

		fmt.Printf("updating %s\n", setting)

		request := client.Request()
		request.
			Method(http.MethodPut).
			Path(payload.Endpoint)

		if payload.Root == "" {
			request.JSON(config)
		} else {
			request.JSON(map[string]interface{}{
				payload.Root: config,
			})
		}

		response, err := request.Send()
		if err != nil {
			return err
		}

		if !response.Ok {
			return fmt.Errorf(
				"could not update %q, the Ops Manager API returned an error:\n%s",
				setting,
				response.Bytes(),
			)
		}
	}

	return nil
}

func (p Payloads) Collect(client *gentleman.Client) (configs, error) {
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

		request := client.Request()
		request.
			Method(http.MethodGet).
			Path(payload.Endpoint)

		response, err := request.Send()
		if err != nil {
			return configs, err
		}

		if !response.Ok {
			return configs, fmt.Errorf(
				"could not gather %q, the Ops Manager API returned an error:\n%s",
				setting,
				response.Bytes(),
			)
		}

		config := map[string]interface{}{}

		err = response.JSON(&config)
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
