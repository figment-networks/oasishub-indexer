package debondingdelegation

import (
	"github.com/figment-networks/oasis-rpc-proxy/grpc/debondingdelegation/debondingdelegationpb"
	"github.com/figment-networks/oasishub-indexer/types"
)

type ListItem struct {
	ValidatorUID string         `json:"validator_uid, omitempty"`
	DelegatorUID string         `json:"delegator_uid"`
	Shares       types.Quantity `json:"shares"`
	DebondEnd    uint64         `json:"debond_end"`
}

type ListView struct {
	Items []ListItem `json:"items"`
}

func ToListView(rawDebondingDelegations map[string]*debondingdelegationpb.DebondingDelegationEntry) *ListView {
	var items []ListItem
	for validatorUID, delegationsMap := range rawDebondingDelegations {
		for delegatorUID, infoArray := range delegationsMap.GetEntries() {
			for _, delegation := range infoArray.GetDebondingDelegations() {
				item := ListItem{
					ValidatorUID: validatorUID,
					DelegatorUID: delegatorUID,
					Shares:       types.NewQuantityFromBytes(delegation.GetShares()),
					DebondEnd:    delegation.GetDebondEndTime(),
				}

				items = append(items, item)
			}
		}
	}
	return &ListView{
		Items: items,
	}
}

func ToListViewForAddress(rawDelegations map[string]*debondingdelegationpb.DebondingDelegationInnerEntry) *ListView {
	var items []ListItem
	for delegatorUID, infoArray := range rawDelegations {
		for _, delegation := range infoArray.GetDebondingDelegations() {
			item := ListItem{
				DelegatorUID: delegatorUID,
				Shares:       types.NewQuantityFromBytes(delegation.GetShares()),
				DebondEnd:    delegation.GetDebondEndTime(),
			}

			items = append(items, item)
		}
	}

	return &ListView{
		Items: items,
	}
}