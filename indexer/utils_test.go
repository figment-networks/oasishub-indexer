package indexer

import (
	"bytes"
	"encoding/binary"
	"math/rand"
	"testing"

	"github.com/figment-networks/oasis-rpc-proxy/grpc/block/blockpb"
	"github.com/figment-networks/oasis-rpc-proxy/grpc/debondingdelegation/debondingdelegationpb"
	"github.com/figment-networks/oasis-rpc-proxy/grpc/delegation/delegationpb"
	"github.com/figment-networks/oasis-rpc-proxy/grpc/state/statepb"
	"github.com/figment-networks/oasis-rpc-proxy/grpc/transaction/transactionpb"
	"github.com/figment-networks/oasis-rpc-proxy/grpc/validator/validatorpb"
	"github.com/golang/protobuf/ptypes"

	"github.com/pkg/errors"
)

var (
	errTestClient   = errors.New("clientErr")
	errTestDbSave   = errors.New("db error on save")
	errTestDbCreate = errors.New("db error on create")
	errTestDbFind   = errors.New("db error on find")

	letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

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

func uintToBytes(num uint64, t *testing.T) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, num)
	if err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

type testBlockOption func(*blockpb.Block)

func setBlockProposerAddress(addr string) testBlockOption {
	return func(b *blockpb.Block) {
		if b.GetHeader() == nil {
			b.Header = &blockpb.Header{}
		}
		b.Header.ProposerAddress = addr
	}
}

func setBlockLastCommitVotes(votes ...*blockpb.Vote) testBlockOption {
	return func(b *blockpb.Block) {
		if b.GetLastCommit() == nil {
			b.LastCommit = &blockpb.Commit{}
		}
		b.LastCommit.Votes = votes
	}
}

func testpbBlock(opts ...testBlockOption) *blockpb.Block {
	b := &blockpb.Block{
		Header: &blockpb.Header{
			ProposerAddress: randString(5),
			ChainId:         randString(5),
			Time:            ptypes.TimestampNow(),
		},
	}

	for _, opt := range opts {
		opt(b)
	}

	return b
}

type testStakingOption func(*statepb.Staking)

func setStakingDelegationEntry(entityID, entryID string, shares []byte) testStakingOption {
	return func(s *statepb.Staking) {
		if s.GetDelegations() == nil {
			s.Delegations = make(map[string]*delegationpb.DelegationEntry)
		}

		if _, ok := s.Delegations[entityID]; !ok {
			s.Delegations[entityID] = &delegationpb.DelegationEntry{
				Entries: make(map[string]*delegationpb.Delegation),
			}
		}

		s.Delegations[entityID].GetEntries()[entryID] = &delegationpb.Delegation{
			Shares: shares,
		}
	}
}

func setDebondingDelegationEntry(entityID, entryID string, shares []byte, endTime uint64) testStakingOption {
	return func(s *statepb.Staking) {
		if s.GetDebondingDelegations() == nil {
			s.DebondingDelegations = make(map[string]*debondingdelegationpb.DebondingDelegationEntry)
		}

		if _, ok := s.DebondingDelegations[entityID]; !ok {
			s.DebondingDelegations[entityID] = &debondingdelegationpb.DebondingDelegationEntry{
				Entries: make(map[string]*debondingdelegationpb.DebondingDelegationInnerEntry),
			}
		}

		if s.DebondingDelegations[entityID].GetEntries()[entryID] == nil {
			s.DebondingDelegations[entityID].GetEntries()[entryID] = &debondingdelegationpb.DebondingDelegationInnerEntry{
				DebondingDelegations: []*debondingdelegationpb.DebondingDelegation{},
			}
		}

		s.DebondingDelegations[entityID].GetEntries()[entryID].DebondingDelegations = append(
			s.DebondingDelegations[entityID].GetEntries()[entryID].GetDebondingDelegations(),
			&debondingdelegationpb.DebondingDelegation{
				Shares:        shares,
				DebondEndTime: endTime,
			},
		)
	}
}

func testpbStaking(opts ...testStakingOption) *statepb.Staking {
	s := &statepb.Staking{
		TotalSupply: randBytes(5),
		CommonPool:  randBytes(5),
		Parameters: &statepb.StakingParameters{
			DebondingInterval:   rand.Uint64(),
			MinDelegationAmount: randBytes(5),
		},
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func setStakingTotalSupply(supply []byte) testStakingOption {
	return func(s *statepb.Staking) {
		s.TotalSupply = supply
	}
}

func setStakingCommonPool(pool []byte) testStakingOption {
	return func(s *statepb.Staking) {
		s.CommonPool = pool
	}
}

func setStakingDebondingInterval(val uint64) testStakingOption {
	return func(s *statepb.Staking) {
		s.GetParameters().DebondingInterval = val
	}
}

func setStakingMinDelegationAmount(amt []byte) testStakingOption {
	return func(s *statepb.Staking) {
		s.GetParameters().MinDelegationAmount = amt
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

type testValidatorOption func(*validatorpb.Validator)

func setValidatorAddress(addr string) testValidatorOption {
	return func(v *validatorpb.Validator) {
		v.Address = addr
	}
}

func setTendermintAddress(addr string) testValidatorOption {
	return func(v *validatorpb.Validator) {
		v.TendermintAddress = addr
	}
}

func setValidatorVotingPower(val int64) testValidatorOption {
	return func(v *validatorpb.Validator) {
		v.VotingPower = val
	}
}

func setValidatorEntityID(id string) testValidatorOption {
	return func(v *validatorpb.Validator) {
		if v.GetNode() == nil {
			v.Node = &validatorpb.Node{}
		}
		v.Node.EntityId = id
	}
}

func testpbValidator(opts ...testValidatorOption) *validatorpb.Validator {
	v := &validatorpb.Validator{
		Address:           randString(5),
		TendermintAddress: randString(5),
		VotingPower:       64,
		Node: &validatorpb.Node{
			EntityId: randString(5),
		},
	}
	for _, opt := range opts {
		opt(v)
	}
	return v
}

func testpbVote(validatorIndex, blockIDFlag int64) *blockpb.Vote {
	return &blockpb.Vote{
		BlockIdFlag:    blockIDFlag,
		ValidatorIndex: validatorIndex,
	}
}
