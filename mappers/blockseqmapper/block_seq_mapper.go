package blockseqmapper

import (
	"github.com/figment-networks/oasishub-indexer/mappers/syncablemapper"
	"github.com/figment-networks/oasishub-indexer/models/blockseq"
	"github.com/figment-networks/oasishub-indexer/models/shared"
	"github.com/figment-networks/oasishub-indexer/models/syncable"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
)

func ToSequence(blockSyncable syncable.Model, validatorsSyncable syncable.Model) (*blockseq.Model, errors.ApplicationError) {
	blockData, err := syncablemapper.UnmarshalBlockData(blockSyncable.Data)
	if err != nil {
		return nil, err
	}
	validatorsData, err := syncablemapper.UnmarshalValidatorsData(validatorsSyncable.Data)
	if err != nil {
		return nil, err
	}

	e := &blockseq.Model{
		Sequence: &shared.Sequence{
			ChainId: blockSyncable.ChainId,
			Height:  blockSyncable.Height,
			Time:    blockSyncable.Time,
		},

		Hash:              types.Hash(blockData.Data.Header.LastBlockID.Hash.String()),
		AppVersion:        int64(blockData.Data.Header.Version.App),
		BlockVersion:      int64(blockData.Data.Header.Version.Block),
		TransactionsCount: types.Count(blockData.Data.Header.NumTxs),
	}

	// Get proposer validator data
	for _, rv := range validatorsData.Data {
		pa := blockData.Data.Header.ProposerAddress.String()

		if pa == rv.Address {
			e.ProposerEntityUID = types.PublicKey(rv.Node.EntityID.String())
		}
	}

	if !e.Valid() {
		return nil, errors.NewErrorFromMessage("block sequence not valid", errors.NotValid)
	}

	return e, nil
}

func ToView(s *blockseq.Model) map[string]interface{} {
	return map[string]interface{}{
		"id":       s.ID,
		"height":   s.Height,
		"time":     s.Time,
		"chain_id": s.ChainId,

		"transactions_count": s.TransactionsCount,
		"app_version":        s.AppVersion,
		"block_version":      s.BlockVersion,
	}
}
