package iterators

import (
	"github.com/figment-networks/oasishub-indexer/types"
	"testing"
)

func Test_HeightIterator(t *testing.T) {
	tests := []struct {
		name   string
		start  types.Height
		end    types.Height
		next   types.Height
		length int64
	}{
		{"1-10", types.Height(1), types.Height(10), types.Height(2), 10},
		{"2-10", types.Height(2), types.Height(10), types.Height(3), 9},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := NewHeightIterator(tt.start, tt.end)
			i.Next()
			nextExp := tt.start.Add(types.Height(1))
			if i.Value() != nextExp {
				t.Errorf("next value for %d should be equal to %d", tt.start, nextExp)
			}
		})
	}
}

func Test_HeightIteratorError(t *testing.T) {
	t.Run("start larger than end", func(t *testing.T) {
		i := NewHeightIterator(10, 1)
		if i.Error() == nil {
			t.Error("start cannot be larger than end")
		}
	})
}