package indexer

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/figment-networks/oasis-rpc-proxy/grpc/validator/validatorpb"
	"github.com/figment-networks/oasishub-indexer/model"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/logger"
	"github.com/pkg/errors"
)

var (
	testClientErr = errors.New("clientErr")

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

func testpbValidator(key string) *validatorpb.Validator {
	return &validatorpb.Validator{
		Address:     randString(5),
		VotingPower: 64,
		Node: &validatorpb.Node{
			EntityId: key,
		},
	}
}
