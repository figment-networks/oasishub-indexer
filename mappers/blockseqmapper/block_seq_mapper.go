package blockseqmapper

import (
	"github.com/figment-networks/oasishub-indexer/db/timescale/orm"
	"github.com/figment-networks/oasishub-indexer/domain/blockdomain"
	"github.com/figment-networks/oasishub-indexer/domain/commons"
	"github.com/figment-networks/oasishub-indexer/domain/syncabledomain"
	"github.com/figment-networks/oasishub-indexer/mappers/syncablemapper"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
)

func FromPersistence(o orm.BlockSeqModel) (*blockdomain.BlockSeq, errors.ApplicationError) {
	e := &blockdomain.BlockSeq{
		DomainEntity: commons.NewDomainEntity(commons.EntityProps{
			ID: o.ID,
		}),
		Sequence: commons.NewSequence(commons.SequenceProps{
			ChainId: o.ChainId,
			Height:  o.Height,
			Time:    o.Time,
		}),

		AppVersion:        o.AppVersion,
		BlockVersion:      o.BlockVersion,
		TransactionsCount: o.TransactionsCount,
	}

	if !e.Valid() {
		return nil, errors.NewErrorFromMessage("block sequence not valid", errors.NotValid)
	}

	return e, nil
}

func ToPersistence(b *blockdomain.BlockSeq) (*orm.BlockSeqModel, errors.ApplicationError) {
	if !b.Valid() {
		return nil, errors.NewErrorFromMessage("block sequence not valid", errors.NotValid)
	}

	return &orm.BlockSeqModel{
		EntityModel: orm.EntityModel{ID: b.ID},
		SequenceModel: orm.SequenceModel{
			ChainId: b.ChainId,
			Height:  b.Height,
			Time:    b.Time,
		},

		AppVersion:        b.AppVersion,
		BlockVersion:      b.BlockVersion,
		TransactionsCount: b.TransactionsCount,
	}, nil
}

func FromData(blockSyncable syncabledomain.Syncable) (*blockdomain.BlockSeq, errors.ApplicationError) {
	blockData, err := syncablemapper.UnmarshalBlockData(blockSyncable.Data)
	if err != nil {
		return nil, err
	}

	e := &blockdomain.BlockSeq{
		DomainEntity: commons.NewDomainEntity(commons.EntityProps{}),
		Sequence: commons.NewSequence(commons.SequenceProps{
			ChainId: blockSyncable.ChainId,
			Height:  blockSyncable.Height,
			Time:    blockSyncable.Time,
		}),

		AppVersion:        int64(blockData.Data.Header.Version.App),
		BlockVersion:      int64(blockData.Data.Header.Version.Block),
		TransactionsCount: types.Count(blockData.Data.Header.NumTxs),
	}

	if !e.Valid() {
		return nil, errors.NewErrorFromMessage("block sequence not valid", errors.NotValid)
	}

	return e, nil
}

func ToView(s *syncabledomain.Syncable) map[string]interface{} {
	return map[string]interface{}{
		"id":        s.ID,
		"height":    s.Height,
		"time":      s.Time,
		"report_id": s.ReportID,
		"chain_id":  s.ChainId,
	}
}
