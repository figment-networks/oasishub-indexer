package indexer

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/figment-networks/oasis-rpc-proxy/grpc/block/blockpb"
	"github.com/figment-networks/oasis-rpc-proxy/grpc/state/statepb"
	"github.com/figment-networks/oasis-rpc-proxy/grpc/transaction/transactionpb"
	"github.com/figment-networks/oasis-rpc-proxy/grpc/validator/validatorpb"
	"github.com/golang/protobuf/ptypes"

	"github.com/figment-networks/oasishub-indexer/model"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/logger"
	"github.com/pkg/errors"
)

var (
	errTestClient = errors.New("clientErr")

	letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func setup(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	err := logger.InitTestLogger()
	if err != nil {
		t.Fatal(err)
	}
}

func ctxWithReport(modelID types.ID) context.Context {
	ctx := context.Background()
	report := &model.Report{
		Model: &model.Model{ID: modelID},
	}

	return context.WithValue(ctx, CtxReport, report)
}

func randString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func randBytes(n int) []byte {
	token := make([]byte, n)
	rand.Read(token)
	return token
}

func testpbBlock() *blockpb.Block {
	return &blockpb.Block{
		Header: &blockpb.Header{
			ProposerAddress: randString(5),
			ChainId:         randString(5),
			Time:            ptypes.TimestampNow(),
		},
	}
}

func testpbStaking() *statepb.Staking {
	return &statepb.Staking{
		TotalSupply: randBytes(5),
		CommonPool:  randBytes(5),
	}
}

func testpbState() *statepb.State {
	return &statepb.State{
		ChainID: randString(5),
		Height:  89,
		Staking: testpbStaking(),
	}
}

func testpbTransaction(key string) *transactionpb.Transaction {
	return &transactionpb.Transaction{
		Hash:      randString(5),
		PublicKey: key,
		Signature: randString(5),
		GasPrice:  randBytes(5),
	}
}

func testpbValidator(key string) *validatorpb.Validator {
	return &validatorpb.Validator{
		Address:     randString(5),
		VotingPower: 64,
		Node: &validatorpb.Node{
			EntityId: key,
		},
	}
}
