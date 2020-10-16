package delegation

import (
	"github.com/figment-networks/oasis-rpc-proxy/grpc/delegation/delegationpb"
	"github.com/figment-networks/oasishub-indexer/types"
)

type ListItem struct {
	ValidatorUID string         `json:"validator_uid, omitempty"`
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

func ToListViewForAddress(rawDelegations map[string]*delegationpb.Delegation) *ListView {
	var items []ListItem
	for delegatorUID, info := range rawDelegations {
		item := ListItem{
			DelegatorUID: delegatorUID,
			Shares:       types.NewQuantityFromBytes(info.GetShares()),
		}

		items = append(items, item)
	}

	return &ListView{
		Items: items,
	}
}