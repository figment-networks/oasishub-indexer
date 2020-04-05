package syncablemapper

import (
	"encoding/json"
	"github.com/figment-networks/oasis-rpc-proxy/controllers"
	"github.com/figment-networks/oasishub-indexer/models/shared"
	"github.com/figment-networks/oasishub-indexer/models/syncable"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
)

func ToSequenceProps(bytes []byte) (*shared.Sequence, errors.ApplicationError) {
	input, err := UnmarshalBlockData(types.Jsonb{RawMessage: bytes})
	if err != nil {
		return nil, err
	}

	return &shared.Sequence{
		ChainId: input.Data.Header.ChainID,
		Height:  types.Height(input.Data.Header.Height),
		Time:    input.Data.Header.Time,
	}, nil
}

func FromProxy(syncableType syncable.Type, sequence shared.Sequence, bytes []byte) (*syncable.Model, errors.ApplicationError) {
	data := types.Jsonb{RawMessage: bytes}

	e := &syncable.Model{
		Sequence: &sequence,
		Data: data,
		Type: syncableType,
	}

	if !e.Valid() {
		return nil, errors.NewErrorFromMessage("syncable not valid", errors.NotValid)
	}
	return e, nil
}

func UnmarshalBlockData(data types.Jsonb) (*controllers.GetBlockResponse, errors.ApplicationError) {
	bytes, err := data.RawMessage.MarshalJSON()
	if err != nil {
		return nil, errors.NewError("some", errors.UnmarshalError, err)
	}
	var input controllers.GetBlockResponse
	if err := json.Unmarshal(bytes, &input); err != nil {
		return nil, errors.NewError("error when trying to unmarshal block response", errors.UnmarshalError, err)
	}
	return &input, nil
}

func UnmarshalStateData(data types.Jsonb) (*controllers.GetConsensusStateResponse, errors.ApplicationError) {
	bytes, err := data.RawMessage.MarshalJSON()
	if err != nil {
		return nil, errors.NewError("some", errors.UnmarshalError, err)
	}
	var input controllers.GetConsensusStateResponse
	if err := json.Unmarshal(bytes, &input); err != nil {
		return nil, errors.NewError("error when trying to unmarshal state response", errors.UnmarshalError, err)
	}
	return &input, nil
}

func UnmarshalValidatorsData(data types.Jsonb) (*controllers.GetValidatorsResponse, errors.ApplicationError) {
	bytes, err := data.RawMessage.MarshalJSON()
	if err != nil {
		return nil, errors.NewError("some", errors.UnmarshalError, err)
	}
	var input controllers.GetValidatorsResponse
	if err := json.Unmarshal(bytes, &input); err != nil {
		return nil, errors.NewError("error when trying to unmarshal validators response", errors.UnmarshalError, err)
	}
	return &input, nil
}

func UnmarshalTransactionsData(data types.Jsonb) (*controllers.GetTransactionsResponse, errors.ApplicationError) {
	bytes, err := data.RawMessage.MarshalJSON()
	if err != nil {
		return nil, errors.NewError("some", errors.UnmarshalError, err)
	}
	var input controllers.GetTransactionsResponse
	if err := json.Unmarshal(bytes, &input); err != nil {
		return nil, errors.NewError("error when trying to unmarshal transactions response", errors.UnmarshalError, err)
	}
	return &input, nil
}