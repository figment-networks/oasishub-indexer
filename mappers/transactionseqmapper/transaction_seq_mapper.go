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

type ListView struct {
	Items []transactionseq.Model `json:"items"`
}

func ToListView(ts []transactionseq.Model) *ListView {
	return &ListView{
		Items: ts,
	}
}
