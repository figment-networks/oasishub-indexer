package shared

import (
	"github.com/figment-networks/oasishub-indexer/clients/grpcclient"
	"github.com/figment-networks/oasishub-indexer/clients/restclient"
	"github.com/figment-networks/oasishub-indexer/clients/timescaleclient"
	"github.com/figment-networks/oasishub-indexer/config"
)

func NewNodeClient() restclient.HttpGetter {
	return restclient.New(restclient.Config{BaseUrl: config.ProxyUrl()})

}

func NewProxyClient() *grpcclient.Client {
	return grpcclient.New(grpcclient.Config{Url: config.ProxyUrl()})
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

