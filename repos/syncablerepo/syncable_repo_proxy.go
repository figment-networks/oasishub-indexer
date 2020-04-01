package syncablerepo

import (
	"fmt"
	"github.com/figment-networks/oasishub-indexer/clients/restclient"
	"github.com/figment-networks/oasishub-indexer/domain/commons"
	"github.com/figment-networks/oasishub-indexer/domain/syncabledomain"
	"github.com/figment-networks/oasishub-indexer/mappers/syncablemapper"
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
	GetHead() (*syncabledomain.Syncable, errors.ApplicationError)
	GetByHeight(syncabledomain.Type, types.Height) (*syncabledomain.Syncable, errors.ApplicationError)
	GetDataByHeight(syncabledomain.Type, types.Height) ([]byte, errors.ApplicationError)
	GetSequencePropsByHeight(types.Height) (*commons.SequenceProps, errors.ApplicationError)
}

type proxyRepo struct {
	rest restclient.HttpGetter
}

func NewProxyRepo(rest restclient.HttpGetter) ProxyRepo {
	return &proxyRepo{
		rest: rest,
	}
}

func (r *proxyRepo) GetHead() (*syncabledomain.Syncable, errors.ApplicationError) {
	return r.GetByHeight(syncabledomain.BlockType, LatestHeight)
}

func (r *proxyRepo) GetByHeight(t syncabledomain.Type, h types.Height) (*syncabledomain.Syncable, errors.ApplicationError) {
	sequenceProps, err := r.GetSequencePropsByHeight(h)
	if err != nil {
		return nil, err
	}

	data, err := r.GetDataByHeight(t, h)
	if err != nil {
		return nil, err
	}

	return syncablemapper.FromData(t, *sequenceProps, data)
}

func (r *proxyRepo) GetDataByHeight(syncableType syncabledomain.Type, height types.Height) ([]byte, errors.ApplicationError) {
	var url string
	switch syncableType {
	case syncabledomain.BlockType:
		url = fmt.Sprintf("block/%d", height)
	case syncabledomain.StateType:
		url = fmt.Sprintf("consensus/%d/state", height)
	case syncabledomain.ValidatorsType:
		url = fmt.Sprintf("validators/%d", height)
	case syncabledomain.TransactionsType:
		url = fmt.Sprintf("transactions/%d", height)
	default:
		return nil, errors.NewErrorFromMessage(fmt.Sprintf("syncable type %s not found", syncableType), errors.ProxyRequestError)
	}
	return r.makeRequest(url)
}

func (r *proxyRepo) GetSequencePropsByHeight(h types.Height) (*commons.SequenceProps, errors.ApplicationError) {
	bytes, err := r.GetDataByHeight(syncabledomain.BlockType, h)
	if err != nil {
		return nil, err
	}
	return syncablemapper.ToSequenceProps(bytes)
}

/*************** Private ***************/

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