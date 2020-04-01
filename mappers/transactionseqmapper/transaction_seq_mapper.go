package transactionseqmapper

import (
	"github.com/figment-networks/oasishub/db/timescale/orm"
	"github.com/figment-networks/oasishub/domain/commons"
	"github.com/figment-networks/oasishub/domain/syncabledomain"
	"github.com/figment-networks/oasishub/domain/transactiondomain"
	"github.com/figment-networks/oasishub/mappers/syncablemapper"
	"github.com/figment-networks/oasishub/types"
	"github.com/figment-networks/oasishub/utils/errors"
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

		Hash:     b.Hash,
		Fee:      b.Fee,
		GasLimit: b.GasLimit,
		GasPrice: b.GasPrice,
		Method:   b.Method,
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

		Hash:     e.Hash,
		Fee:      e.Fee,
		GasLimit: e.GasLimit,
		GasPrice: e.GasPrice,
		Method:   e.Method,
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

			Hash:     types.Hash(rv.Hash),
			Fee:      rv.Fee.Int64(),
			GasLimit: rv.GasLimit,
			GasPrice: rv.GasPrice.Int64(),
			Method:   rv.Method,
		}

		if !e.Valid() {
			return nil, errors.NewErrorFromMessage("transaction sequence not valid", errors.NotValid)
		}

		transactions = append(transactions, e)
	}
	return transactions, nil
}
