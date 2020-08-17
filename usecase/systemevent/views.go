package systemevent

import (
	"github.com/figment-networks/oasishub-indexer/model"
	"github.com/figment-networks/oasishub-indexer/types"
)

type ListItem struct {
	*model.Model

	Height int64       `json:"height"`
	Time   types.Time  `json:"time"`
	Actor  string      `json:"actor"`
	Kind   string      `json:"kind"`
	Data   types.Jsonb `json:"data"`
}

type ListView struct {
	Items []ListItem `json:"items"`
}

func ToListView(validators []model.SystemEvent) *ListView {
	var items []ListItem
	for _, m := range validators {
		item := ListItem{
			Actor:  m.Actor,
			Height: m.Height,
			Time:   m.Time,
			Kind:   m.Kind.String(),
			Data:   m.Data,
		}

		items = append(items, item)
	}

	return &ListView{
		Items: items,
	}
}
