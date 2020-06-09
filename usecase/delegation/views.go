package delegation

import (
	"github.com/figment-networks/oasis-rpc-proxy/grpc/delegation/delegationpb"
	"github.com/figment-networks/oasishub-indexer/types"
)

type ListItem struct {
	ValidatorUID string         `json:"validator_uid"`
	DelegatorUID string         `json:"delegator_uid"`
	Shares       types.Quantity `json:"shares"`
}

type ListView struct {
	Items []ListItem `json:"items"`
}

func ToListView(rawDelegations map[string]*delegationpb.DelegationEntry) *ListView {
	var items []ListItem
	for validatorUID, delegationsMap := range rawDelegations {
		for delegatorUID, info := range delegationsMap.GetEntries() {
			item := ListItem{
				ValidatorUID: validatorUID,
				DelegatorUID: delegatorUID,
				Shares:       types.NewQuantityFromBytes(info.GetShares()),
			}

			items = append(items, item)
		}
	}

	return &ListView{
		Items: items,
	}
}
