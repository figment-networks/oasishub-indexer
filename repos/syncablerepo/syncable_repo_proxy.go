package syncablerepo

import (
	"context"
	"fmt"
	"github.com/figment-networks/oasis-rpc-proxy/grpc/block/blockpb"
	"github.com/figment-networks/oasis-rpc-proxy/grpc/state/statepb"
	"github.com/figment-networks/oasis-rpc-proxy/grpc/transaction/transactionpb"
	"github.com/figment-networks/oasis-rpc-proxy/grpc/validator/validatorpb"
	"github.com/figment-networks/oasishub-indexer/mappers/syncablemapper"
	"github.com/figment-networks/oasishub-indexer/models/shared"
	"github.com/figment-networks/oasishub-indexer/models/syncable"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
	"github.com/golang/protobuf/proto"
	"sync"
)

const (
	LatestHeight = 0
)

type ProxyRepo interface {
	//Queries
	GetHead() (*syncable.Model, errors.ApplicationError)
	GetByHeight(syncable.Type, types.Height) (*syncable.Model, errors.ApplicationError)
}

type sequencePropsCache struct {
	items map[types.Height]*shared.Sequence

	mu sync.Mutex
}

type proxyRepo struct {
	blockClient       blockpb.BlockServiceClient
	stateClient       statepb.StateServiceClient
	transactionClient transactionpb.TransactionServiceClient
	validatorClient   validatorpb.ValidatorServiceClient

	sequencePropsCache sequencePropsCache
}

func NewProxyRepo(
	blockClient blockpb.BlockServiceClient,
	stateClient statepb.StateServiceClient,
	transactionClient transactionpb.TransactionServiceClient,
	validatorClient validatorpb.ValidatorServiceClient,
) ProxyRepo {
	return &proxyRepo{
		blockClient:       blockClient,
		stateClient:       stateClient,
		transactionClient: transactionClient,
		validatorClient:   validatorClient,

		sequencePropsCache: sequencePropsCache{
			items: map[types.Height]*shared.Sequence{},
			mu:    sync.Mutex{},
		},
	}
}

func (r *proxyRepo) GetHead() (*syncable.Model, errors.ApplicationError) {
	return r.GetByHeight(syncable.BlockType, LatestHeight)
}

func (r *proxyRepo) GetByHeight(t syncable.Type, h types.Height) (*syncable.Model, errors.ApplicationError) {
	ctx := context.Background()
	sequenceProps, err := r.getSequencePropsByHeight(ctx, h)
	if err != nil {
		return nil, err
	}

	data, err := r.getRawDataByHeight(ctx, t, h)
	if err != nil {
		return nil, err
	}

	return syncablemapper.FromProxy(t, *sequenceProps, data)
}

/*************** Private ***************/

func (r *proxyRepo) getSequencePropsByHeight(ctx context.Context, h types.Height) (*shared.Sequence, errors.ApplicationError) {
	r.sequencePropsCache.mu.Lock()
	defer r.sequencePropsCache.mu.Unlock()
	sequenceProps, ok := r.sequencePropsCache.items[h]
	if !ok {
		data, err := r.getRawDataByHeight(ctx, syncable.BlockType, h)
		if err != nil {
			return nil, err
		}

		res := data.(*blockpb.GetByHeightResponse)

		sequenceProps = &shared.Sequence{
			ChainId: res.Block.Header.GetChainId(),
			Height:  types.Height(res.Block.Header.GetHeight()),
			Time:    *types.NewTimeFromTimestamp(*res.Block.Header.GetTime()),
		}
		r.sequencePropsCache.items[h] = sequenceProps
	}

	return sequenceProps, nil
}

func (r *proxyRepo) getRawDataByHeight(ctx context.Context, syncableType syncable.Type, h types.Height) (proto.Message, errors.ApplicationError) {
	var res proto.Message
	var err error
	switch syncableType {
	case syncable.BlockType:
		res, err = r.blockClient.GetByHeight(ctx, &blockpb.GetByHeightRequest{Height: h.Int64()})
		if err != nil {
			return nil, errors.NewError("error getting block by height", errors.ProxyRequestError, err)
		}
	case syncable.StateType:
		res, err = r.stateClient.GetByHeight(ctx, &statepb.GetByHeightRequest{Height: h.Int64()})
		if err != nil {
			return nil, errors.NewError("error getting state by height", errors.ProxyRequestError, err)
		}
	case syncable.ValidatorsType:
		res, err = r.validatorClient.GetByHeight(ctx, &validatorpb.GetByHeightRequest{Height: h.Int64()})
		if err != nil {
			return nil, errors.NewError("error getting validators by height", errors.ProxyRequestError, err)
		}
	case syncable.TransactionsType:
		res, err = r.transactionClient.GetByHeight(ctx, &transactionpb.GetByHeightRequest{Height: h.Int64()})
		if err != nil {
			return nil, errors.NewError("error getting transactions by height", errors.ProxyRequestError, err)
		}
	default:
		return nil, errors.NewErrorFromMessage(fmt.Sprintf("syncable type %s not found", syncableType), errors.ProxyRequestError)
	}
	return res, nil
}
