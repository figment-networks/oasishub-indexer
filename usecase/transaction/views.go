package transaction

import (
	"github.com/figment-networks/oasis-rpc-proxy/grpc/transaction/transactionpb"
	"github.com/figment-networks/oasishub-indexer/types"
)

type ListItem struct {
	PublicKey string         `json:"public_key"`
	Hash      string         `json:"hash"`
	Nonce     uint64         `json:"nonce"`
	Fee       types.Quantity `json:"fee"`
	GasLimit  uint64         `json:"gas_limit"`
	GasPrice  types.Quantity `json:"gas_price"`
	Method    string         `json:"method"`
}

type ListView struct {
	Items []ListItem `json:"items"`
}

func ToListView(rawTransactions []*transactionpb.Transaction) *ListView {
	var items []ListItem
	for _, rawTransaction := range rawTransactions {
		item := ListItem{
			PublicKey: rawTransaction.GetPublicKey(),
			Hash:      rawTransaction.GetHash(),
			Nonce:     rawTransaction.GetNonce(),
			Fee:       types.NewQuantityFromBytes(rawTransaction.GetFee()),
			GasLimit:  rawTransaction.GetGasLimit(),
			GasPrice:  types.NewQuantityFromBytes(rawTransaction.GetGasPrice()),
			Method:    rawTransaction.GetMethod(),
		}

		items = append(items, item)
	}

	return &ListView{
		Items: items,
	}
}
