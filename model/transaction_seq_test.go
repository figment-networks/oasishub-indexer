package model

import (
	"github.com/figment-networks/oasishub-indexer/models/shared"
	"github.com/figment-networks/oasishub-indexer/types"
	"testing"
	"time"
)

func Test_TransactionSeq(t *testing.T) {
	chainId := "chain123"
	model := &shared.Model{}
	seq := &shared.Sequence{
		ChainId: chainId,
		Height:  int64(10),
		Time:    *types.NewTimeFromTime(time.Now()),
	}

	t.Run("validation failed", func(t *testing.T) {
		m := Model{
			Model: model,
			Sequence: seq,
		}

		if m.Valid() {
			t.Errorf("model should not be valid %+v", m)
		}
	})

	t.Run("validation success", func(t *testing.T) {
		m := Model{
			Model: model,
			Sequence: seq,

			PublicKey: string("pb-test"),
			Hash: string("hash-test"),
			Nonce: uint64(10),
		}

		if !m.Valid() {
			t.Errorf("model should be valid %+v", m)
		}
	})

	t.Run("models not equal", func(t *testing.T) {
		m1 := Model{
			Model: model,
			Sequence: seq,

			PublicKey: string("pb-test"),
			Hash: string("hash-test"),
			Nonce: uint64(10),
		}
		m2 := Model{
			Model: model,
			Sequence: seq,

			PublicKey: string("pb-test-2"),
			Hash: string("hash-test-2"),
			Nonce: uint64(20),
		}

		if m1.Equal(m2) {
			t.Errorf("models should not be equal first: %+v; second %+v", m1, m2)
		}
	})

	t.Run("models equal", func(t *testing.T) {
		m1 := Model{
			Model: model,
			Sequence: seq,

			PublicKey: string("pb-test"),
			Hash: string("hash-test"),
			Nonce: uint64(10),
		}
		m2 := Model{
			Model: model,
			Sequence: seq,

			PublicKey: string("pb-test"),
			Hash: string("hash-test"),
			Nonce: uint64(10),
		}

		if !m1.Equal(m2) {
			t.Errorf("models should be equal first: %+v; second %+v", m1, m2)
		}
	})
}


