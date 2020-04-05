package transactionseqmapper

import (
	"github.com/figment-networks/oasishub-indexer/mappers/syncablemapper"
	"github.com/figment-networks/oasishub-indexer/models/shared"
	"github.com/figment-networks/oasishub-indexer/models/syncable"
	"github.com/figment-networks/oasishub-indexer/models/transactionseq"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
)

func ToSequence(transactionsSyncable syncable.Model) ([]transactionseq.Model, errors.ApplicationError) {
	transactionsData, err := syncablemapper.UnmarshalTransactionsData(transactionsSyncable.Data)
	if err != nil {
		return nil, err
	}

	var transactions []transactionseq.Model
	for _, rv := range transactionsData.Data {
		e := transactionseq.Model{
			Sequence: &shared.Sequence{
				ChainId: transactionsSyncable.ChainId,
				Height:  transactionsSyncable.Height,
				Time:    transactionsSyncable.Time,
			},

			PublicKey: types.PublicKey(rv.PublicKey),
			Hash:      types.Hash(rv.Hash),
			Nonce:     types.Nonce(rv.Nonce),
			Fee:       rv.Fee.Int64(),
			GasLimit:  rv.GasLimit,
			GasPrice:  rv.GasPrice.Int64(),
			Method:    rv.Method,
		}

		if !e.Valid() {
			return nil, errors.NewErrorFromMessage("transaction sequence not valid", errors.NotValid)
		}

		transactions = append(transactions, e)
	}
	return transactions, nil
}

func ToView(ts []*transactionseq.Model) []map[string]interface{} {
	var items []map[string]interface{}
	for _, t := range ts {
		i := map[string]interface{}{
			"id":         t.ID,
			"height":     t.Height,
			"time":       t.Time,
			"chain_id":   t.ChainId,

			"public_key": t.PublicKey,
			"hash":       t.Hash,
			"nonce":      t.Nonce,
			"gas_price":  t.GasPrice,
			"gas_limit":  t.GasLimit,
			"fee":        t.Fee,
			"method":     t.Method,
		}
		items = append(items, i)
	}
	return items
}
