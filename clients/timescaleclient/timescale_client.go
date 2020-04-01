package timescaleclient

import (
	"fmt"
	"github.com/figment-networks/oasishub-indexer/config"
	"github.com/figment-networks/oasishub-indexer/utils/log"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

type DbClient interface{
	Client() *gorm.DB
}

type tsClient struct {
	dsn string
	c *gorm.DB
}

var _ DbClient = (*tsClient)(nil)

func New(props Config) *tsClient {
	props.isValid()
	dsn := fmt.Sprintf("user=%s password=%s host=%s dbname=%s sslmode=%s", props.User, props.Password, props.Host, props.DatabaseName, props.SLLMode)

	log.Info("initializing data source...", log.Field("type", "database"))

	c, err := gorm.Open("postgres", dsn)
	if err != nil {
		log.Error(err)
		panic("could not connect to data source")
	}

	log.Info("data source initialized successfully", log.Field("type", "database"))

	c.LogMode(config.DbDetailedLog())

	registerPlugins(c)

	return &tsClient{
		dsn: dsn,
		c: c,
	}
}

func (c *tsClient) Client() *gorm.DB {
	return c.c
}
