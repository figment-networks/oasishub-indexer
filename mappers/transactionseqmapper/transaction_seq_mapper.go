package transactionseqmapper

import (
	"github.com/figment-networks/oasishub-indexer/db/timescale/orm"
	"github.com/figment-networks/oasishub-indexer/domain/commons"
	"github.com/figment-networks/oasishub-indexer/domain/syncabledomain"
	"github.com/figment-networks/oasishub-indexer/domain/transactiondomain"
	"github.com/figment-networks/oasishub-indexer/mappers/syncablemapper"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
)

func FromPersistence(b orm.TransactionSeqModel) (*transactiondomain.TransactionSeq, errors.ApplicationError) {
	e := &transactiondomain.TransactionSeq{
		DomainEntity: commons.NewDomainEntity(commons.EntityProps{
			ID: b.ID,
		}),
		Sequence: commons.NewSequence(commons.SequenceProps{
			ChainId: b.ChainId,
			Height:  b.Height,
			Time:    b.Time,
		}),

		PublicKey: b.PublicKey,
		Hash:      b.Hash,
		Nonce:     b.Nonce,
		Fee:       b.Fee,
		GasLimit:  b.GasLimit,
		GasPrice:  b.GasPrice,
		Method:    b.Method,
	}

	if !e.Valid() {
		return nil, errors.NewErrorFromMessage("transaction sequence not valid", errors.NotValid)
	}

	return e, nil
}

func ToPersistence(e *transactiondomain.TransactionSeq) (*orm.TransactionSeqModel, errors.ApplicationError) {
	if !e.Valid() {
		return nil, errors.NewErrorFromMessage("transaction sequence not valid", errors.NotValid)
	}

	return &orm.TransactionSeqModel{
		EntityModel: orm.EntityModel{ID: e.ID},
		SequenceModel: orm.SequenceModel{
			ChainId: e.ChainId,
			Height:  e.Height,
			Time:    e.Time,
		},

		PublicKey: e.PublicKey,
		Hash:      e.Hash,
		Nonce:     e.Nonce,
		Fee:       e.Fee,
		GasLimit:  e.GasLimit,
		GasPrice:  e.GasPrice,
		Method:    e.Method,
	}, nil
}

func FromData(transactionsSyncable syncabledomain.Syncable) ([]*transactiondomain.TransactionSeq, errors.ApplicationError) {
	transactionsData, err := syncablemapper.UnmarshalTransactionsData(transactionsSyncable.Data)
	if err != nil {
		return nil, err
	}

	var transactions []*transactiondomain.TransactionSeq
	for _, rv := range transactionsData.Data {
		e := &transactiondomain.TransactionSeq{
			DomainEntity: commons.NewDomainEntity(commons.EntityProps{}),
			Sequence: commons.NewSequence(commons.SequenceProps{
				ChainId: transactionsSyncable.ChainId,
				Height:  transactionsSyncable.Height,
				Time:    transactionsSyncable.Time,
			}),

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

func ToView(ts []*transactiondomain.TransactionSeq) []map[string]interface{} {
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
