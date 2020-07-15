package indexer

import (
	"context"
	"testing"

	"github.com/figment-networks/oasis-rpc-proxy/grpc/block/blockpb"
	"github.com/figment-networks/oasis-rpc-proxy/grpc/transaction/transactionpb"
	"github.com/figment-networks/oasis-rpc-proxy/grpc/validator/validatorpb"
)

func TestParserTask_Run(t *testing.T) {
	setup(t)
	testProposerAddr := "testProposerAddr"
	testProposerKey := "testProposerKey"

	testblockProposer := testpbValidator(testProposerKey)
	testblockProposer.Address = testProposerAddr

	testblock := testpbBlock()
	testblock.Header.ProposerAddress = testProposerAddr

	tests := []struct {
		description              string
		rawBlock                 *blockpb.Block
		rawTransactions          []*transactionpb.Transaction
		rawValidators            []*validatorpb.Validator
		expectedTransactionCount int64
		expectedUID              string
	}{
		{"updates empty state",
			nil,
			[]*transactionpb.Transaction{},
			[]*validatorpb.Validator{},
			0,
			"",
		},
		{"updates payload.ParsedBlockData.TransactionsCount",
			nil,
			[]*transactionpb.Transaction{testpbTransaction("t1"), testpbTransaction("t2")},
			[]*validatorpb.Validator{},
			2,
			"",
		},
		{"updates payload.ParsedBlockData.ProposerEntityUID when proposer is in validator list",
			testblock,
			[]*transactionpb.Transaction{},
			[]*validatorpb.Validator{testpbValidator("t1"), testblockProposer, testpbValidator("t2")},
			0,
			testProposerKey,
		},
		{"does not update payload.ParsedBlockData.ProposerEntityUID when proposer is not in validator list",
			testblock,
			[]*transactionpb.Transaction{},
			[]*validatorpb.Validator{testpbValidator("t1"), testpbValidator("t2")},
			0,
			"",
		},
		{"does not update payload.ParsedBlockData.ProposerEntityUID when there's no block",
			nil,
			[]*transactionpb.Transaction{},
			[]*validatorpb.Validator{testpbValidator("t1"), testblockProposer, testpbValidator("t2")},
			0,
			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			ctx := context.Background()

			task := NewBlockParserTask()

			pl := &payload{
				RawBlock:        tt.rawBlock,
				RawValidators:   tt.rawValidators,
				RawTransactions: tt.rawTransactions,
			}

			if err := task.Run(ctx, pl); err != nil {
				t.Errorf("unexpected error on Run, want %v; got %v", nil, err)
				return
			}

			if pl.ParsedBlock.TransactionsCount != tt.expectedTransactionCount {
				t.Errorf("Unexpected ProposerEntityUID, want: %+v, got: %+v", tt.expectedTransactionCount, pl.ParsedBlock.TransactionsCount)
				return
			}

			if pl.ParsedBlock.ProposerEntityUID != tt.expectedUID {
				t.Errorf("Unexpected ProposerEntityUID count, want: %+v, got: %+v", tt.expectedUID, pl.ParsedBlock.ProposerEntityUID)
				return
			}
		})
	}
}
