package syncablerepo

import (
	"fmt"
	"github.com/figment-networks/oasishub-indexer/clients/restclient"
	"github.com/figment-networks/oasishub-indexer/mappers/syncablemapper"
	"github.com/figment-networks/oasishub-indexer/models/shared"
	"github.com/figment-networks/oasishub-indexer/models/syncable"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
	"github.com/figment-networks/oasishub-indexer/utils/log"
	"io/ioutil"
)

const (
	LatestHeight = 0
)

type ProxyRepo interface {
	//Queries
	GetHead() (*syncable.Model, errors.ApplicationError)
	GetByHeight(syncable.Type, types.Height) (*syncable.Model, errors.ApplicationError)
}

type proxyRepo struct {
	rest restclient.HttpGetter
}

func NewProxyRepo(rest restclient.HttpGetter) ProxyRepo {
	return &proxyRepo{
		rest: rest,
	}
}

func (r *proxyRepo) GetHead() (*syncable.Model, errors.ApplicationError) {
	return r.GetByHeight(syncable.BlockType, LatestHeight)
}

func (r *proxyRepo) GetByHeight(t syncable.Type, h types.Height) (*syncable.Model, errors.ApplicationError) {
	sequenceProps, err := r.getSequencePropsByHeight(h)
	if err != nil {
		return nil, err
	}

	bytes, err := r.getRawDataByHeight(t, h)
	if err != nil {
		return nil, err
	}

	return syncablemapper.FromProxy(t, *sequenceProps, bytes)
}

/*************** Private ***************/

func (r *proxyRepo) getSequencePropsByHeight(h types.Height) (*shared.Sequence, errors.ApplicationError) {
	//TODO: Can be cached
	bytes, err := r.getRawDataByHeight(syncable.BlockType, h)
	if err != nil {
		return nil, err
	}
	return syncablemapper.ToSequenceProps(bytes)
}

func (r *proxyRepo) getRawDataByHeight(syncableType syncable.Type, height types.Height) ([]byte, errors.ApplicationError) {
	var url string
	switch syncableType {
	case syncable.BlockType:
		url = fmt.Sprintf("block/%d", height)
	case syncable.StateType:
		url = fmt.Sprintf("consensus/%d/state", height)
	case syncable.ValidatorsType:
		url = fmt.Sprintf("validators/%d", height)
	case syncable.TransactionsType:
		url = fmt.Sprintf("transactions/%d", height)
	default:
		return nil, errors.NewErrorFromMessage(fmt.Sprintf("syncable type %s not found", syncableType), errors.ProxyRequestError)
	}
	return r.makeRequest(url)
}

func (r *proxyRepo) makeRequest(url string) ([]byte, errors.ApplicationError) {
	log.Info(fmt.Sprintf("making request to node to get syncable at %s", url), log.Field("type", "proxy"))

	response, err := r.rest.Get(url, nil)
	if err != nil {
		return nil, errors.NewError(fmt.Sprintf("error getting syncable from node at %s", url), errors.ProxyRequestError, err)
	}

	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.NewError("invalid response body from node", errors.ProxyInvalidResponseError, err)
	}
	defer response.Body.Close()

	return bytes, nil
}