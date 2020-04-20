package chainmapper

import (
	"github.com/figment-networks/oasis-rpc-proxy/grpc/chain/chainpb"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/golang/protobuf/ptypes/timestamp"
	"testing"
	"time"
)

func Test_ChainMapper(t *testing.T) {
	t.Run("FromProxy() maps data correctly", func(t *testing.T) {
		tm := time.Now()
		seconds := tm.Unix()
		nanos := int32(tm.Nanosecond())
		res := chainpb.GetCurrentResponse{
			Chain: &chainpb.Chain{
				Id:                   "123456",
				Name:                 "test-chain",
				GenesisTime:          &timestamp.Timestamp{Seconds: seconds, Nanos: nanos},
				Height:               1,
			},
		}

		chain := FromProxy(res)

		if chain.Id != "123456" {
			t.Errorf("id expected: %s, got: %s", "123456", chain.Id)
		}

		if chain.Name != "test-chain" {
			t.Errorf("name expected: %s, got: %s", "test-chain", chain.Name)
		}

		if !chain.GenesisTime.Equal(*types.NewTimeFromTime(tm)) {
			t.Errorf("genesis time expected: %s, got: %s", *types.NewTimeFromTime(tm), chain.GenesisTime)
		}

		if chain.Height != 1 {
			t.Errorf("height expected: %d, got: %d", 1, chain.Height)
		}
	})
}
