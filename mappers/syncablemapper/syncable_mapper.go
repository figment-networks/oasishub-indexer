package syncablemapper

import (
	"fmt"
	"github.com/figment-networks/oasis-rpc-proxy/grpc/block/blockpb"
	"github.com/figment-networks/oasis-rpc-proxy/grpc/state/statepb"
	"github.com/figment-networks/oasis-rpc-proxy/grpc/transaction/transactionpb"
	"github.com/figment-networks/oasis-rpc-proxy/grpc/validator/validatorpb"
	"github.com/figment-networks/oasishub-indexer/models/shared"
	"github.com/figment-networks/oasishub-indexer/models/syncable"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

func FromProxy(syncableType syncable.Type, sequence shared.Sequence, data proto.Message) (*syncable.Model, errors.ApplicationError) {
	var bytes string
	var err error
	marshaler := jsonpb.Marshaler{}

	switch syncableType {
	case syncable.BlockType:
		res := data.(*blockpb.GetByHeightResponse)
		bytes, err = marshaler.MarshalToString(res)
	case syncable.StateType:
		res := data.(*statepb.GetByHeightResponse)
		bytes, err = marshaler.MarshalToString(res)
	case syncable.ValidatorsType:
		res := data.(*validatorpb.GetByHeightResponse)
		bytes, err = marshaler.MarshalToString(res)
	case syncable.TransactionsType:
		res := data.(*transactionpb.GetByHeightResponse)
		bytes, err = marshaler.MarshalToString(res)
	default:
		return nil, errors.NewErrorFromMessage(fmt.Sprintf("syncable type %s not found", syncableType), errors.ProxyRequestError)
	}

	if err != nil {
		return nil, errors.NewErrorFromMessage(fmt.Sprintf("syncable type %s could not be marshaled to JSON", syncableType), errors.ProxyUnmarshalError)
	}

	d := types.Jsonb{RawMessage: []byte(bytes)}

	e := &syncable.Model{
		Sequence: &sequence,
		Data: d,
		Type: syncableType,
	}

	if !e.Valid() {
		return nil, errors.NewErrorFromMessage("syncable not valid", errors.NotValid)
	}
	return e, nil
}

func UnmarshalBlockData(data types.Jsonb) (*blockpb.Block, errors.ApplicationError) {
	res := &blockpb.GetByHeightResponse{}
	err := jsonpb.UnmarshalString(string(data.RawMessage), res)
	if err != nil {
		return nil, errors.NewError("could not unmarshal grpc block response", errors.ProxyUnmarshalError, err)
	}

	return res.Block, nil
}

func UnmarshalStateData(data types.Jsonb) (*statepb.State, errors.ApplicationError) {
	res := &statepb.GetByHeightResponse{}
	err := jsonpb.UnmarshalString(string(data.RawMessage), res)
	if err != nil {
		return nil, errors.NewError("could not unmarshal grpc state response", errors.ProxyUnmarshalError, err)
	}

	return res.State, nil
}

func UnmarshalValidatorsData(data types.Jsonb) ([]*validatorpb.Validator, errors.ApplicationError) {
	res := &validatorpb.GetByHeightResponse{}
	err := jsonpb.UnmarshalString(string(data.RawMessage), res)
	if err != nil {
		return nil, errors.NewError("could not unmarshal grpc validator response", errors.ProxyUnmarshalError, err)
	}

	return res.Validators, nil
}

func UnmarshalTransactionsData(data types.Jsonb) ([]*transactionpb.Transaction, errors.ApplicationError) {
	res := &transactionpb.GetByHeightResponse{}
	err := jsonpb.UnmarshalString(string(data.RawMessage), res)
	if err != nil {
		return nil, errors.NewError("could not unmarshal grpc transaction response", errors.ProxyUnmarshalError, err)
	}

	return res.Transactions, nil
}