package ossindexapi

import (
	"github.com/a-grasso/deprec/configuration"
	"github.com/a-grasso/deprec/logging"
	"github.com/nscuro/ossindex-client"
)

type Client struct {
	*ossindex.Client
}

func NewClient(config configuration.OSSIndex) *Client {

	client, err := ossindex.NewClient(ossindex.WithAuthentication(config.Username, config.Token))
	if err != nil {
		logging.Logger.Warn("error creating ossindex api client")
		return nil
	}

	return &Client{
		client,
	}
}
