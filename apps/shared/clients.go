package shared

import (
	"github.com/figment-networks/oasishub/clients/restclient"
	"github.com/figment-networks/oasishub/clients/timescaleclient"
	"github.com/figment-networks/oasishub/config"
)

func NewNodeClient() restclient.HttpGetter {
	return restclient.New(restclient.Config{BaseUrl: config.NodeUrl()})

}

func NewDbClient() timescaleclient.DbClient {
	return timescaleclient.New(timescaleclient.Config{
		Host: config.DbHost(),
		User: config.DbUser(),
		Password: config.DbPassword(),
		DatabaseName: config.DbName(),
		SLLMode:      config.DbSSLMode(),
	})
}

