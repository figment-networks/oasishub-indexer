package apr

import (
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCalculateApr(t *testing.T) {
	tests := []struct {
		description         string
		escrowActiveBalance int64
		totalRewards        int64
		expectedRes         float64
		expectedErr         bool
	}{
		{description: "zero total reward",
			escrowActiveBalance: 0,
			totalRewards:        0,
			expectedErr:         true,
		},
		{description: "zero total reward",
			escrowActiveBalance: 102,
			totalRewards:        2,
			expectedRes:         24.33333333333333,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			b := types.NewQuantityFromInt64(tt.escrowActiveBalance)
			r := types.NewQuantityFromInt64(tt.totalRewards)
			apr, err := calculateAPR(b, r)

			if err != nil {
				assert.Equal(t, true, tt.expectedErr)
			} else {
				assert.Equal(t, false, tt.expectedErr)
				assert.Equal(t, apr, tt.expectedRes)
			}
		})
	}

}
