package syncablemapper

import (
	"encoding/json"
	"github.com/figment-networks/oasis-rpc-proxy/controllers"
	"github.com/figment-networks/oasishub-indexer/db/timescale/orm"
	"github.com/figment-networks/oasishub-indexer/domain/commons"
	"github.com/figment-networks/oasishub-indexer/domain/syncabledomain"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
	"github.com/figment-networks/oasishub-indexer/utils/log"
	"github.com/jinzhu/gorm/dialects/postgres"
)

func FromPersistence(b orm.SyncableModel) (*syncabledomain.Syncable, errors.ApplicationError) {
	data, err := b.Data.MarshalJSON()
	if err != nil {
		log.Error(err)
		return nil, errors.NewError("error when trying to marshal syncable data", errors.UnmarshalError, err)
	}

	e := &syncabledomain.Syncable{
		DomainEntity: commons.NewDomainEntity(commons.EntityProps{
			ID: b.ID,
		}),
		Sequence: commons.NewSequence(commons.SequenceProps{
			ChainId:      b.ChainId,
			Height:       b.Height,
			Time:         b.Time,
		}),

		Type:        b.Type,
		ReportID:    b.ReportID,
		Data:        data,
		ProcessedAt: b.ProcessedAt,
	}

	if !e.Valid() {
		return nil, errors.NewErrorFromMessage("syncable not valid", errors.NotValid)
	}

	return e, nil
}

func ToPersistence(b *syncabledomain.Syncable) (*orm.SyncableModel, errors.ApplicationError) {
	if !b.Valid() {
		return nil, errors.NewErrorFromMessage("syncable not valid", errors.NotValid)
	}

	data := postgres.Jsonb{RawMessage: b.Data}

	return &orm.SyncableModel{
		EntityModel: orm.EntityModel{ID: b.ID},
		SequenceModel: orm.SequenceModel{
			ChainId:      b.ChainId,
			Height:       b.Height,
			Time:         b.Time,
		},
		Type:        b.Type,
		ReportID:    b.ReportID,
		Data:        data,
		ProcessedAt: b.ProcessedAt,
	}, nil
}

func ToSequenceProps(bytes []byte) (*commons.SequenceProps, errors.ApplicationError) {
	input, err := UnmarshalBlockData(bytes)
	if err != nil {
		return nil, err
	}

	return &commons.SequenceProps{
		ChainId: input.Data.Header.ChainID,
		Height:  types.Height(input.Data.Header.Height),
		Time:    input.Data.Header.Time,
	}, nil
}

func FromData(syncableType syncabledomain.Type, sequenceProps commons.SequenceProps, data []byte) (*syncabledomain.Syncable, errors.ApplicationError) {
	e := &syncabledomain.Syncable{
		DomainEntity: commons.NewDomainEntity(commons.EntityProps{}),
		Sequence: commons.NewSequence(sequenceProps),
		Data: data,
		Type: syncableType,
	}

	if !e.Valid() {
		return nil, errors.NewErrorFromMessage("syncable not valid", errors.NotValid)
	}

	return e, nil
}

func UnmarshalBlockData(bytes []byte) (*controllers.GetBlockResponse, errors.ApplicationError) {
	var input controllers.GetBlockResponse
	if err := json.Unmarshal(bytes, &input); err != nil {
		log.Error(err)
		return nil, errors.NewError("error when trying to unmarshal block response", errors.UnmarshalError, err)
	}

	return &input, nil
}

func UnmarshalStateData(bytes []byte) (*controllers.GetConsensusStateResponse, errors.ApplicationError) {
	var input controllers.GetConsensusStateResponse
	if err := json.Unmarshal(bytes, &input); err != nil {
		log.Error(err)
		return nil, errors.NewError("error when trying to unmarshal state response", errors.UnmarshalError, err)
	}

	return &input, nil
}

func UnmarshalValidatorsData(bytes []byte) (*controllers.GetValidatorsResponse, errors.ApplicationError) {
	var input controllers.GetValidatorsResponse
	if err := json.Unmarshal(bytes, &input); err != nil {
		log.Error(err)
		return nil, errors.NewError("error when trying to unmarshal validators response", errors.UnmarshalError, err)
	}

	return &input, nil
}

func UnmarshalTransactionsData(bytes []byte) (*controllers.GetTransactionsResponse, errors.ApplicationError) {
	var input controllers.GetTransactionsResponse
	if err := json.Unmarshal(bytes, &input); err != nil {
		log.Error(err)
		return nil, errors.NewError("error when trying to unmarshal transactions response", errors.UnmarshalError, err)
	}

	return &input, nil
}