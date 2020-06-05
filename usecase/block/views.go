package block

import (
	"github.com/figment-networks/oasis-rpc-proxy/grpc/block/blockpb"
	"github.com/figment-networks/oasishub-indexer/types"
)

type DetailsView struct {
	AppVersion         uint64     `json:"app_version"`
	BlockVersion       uint64     `json:"block_version"`
	ChainId            string     `json:"chain_id"`
	Height             int64      `json:"height"`
	Time               types.Time `json:"time"`
	LastBlockIdHash    string     `json:"last_block_id_hash"`
	LastCommitHash     string     `json:"last_commit_hash"`
	DataHash           string     `json:"data_hash"`
	ValidatorsHash     string     `json:"validators_hash"`
	NextValidatorsHash string     `json:"next_validators_hash"`
	ConsensusHash      string     `json:"consensus_hash"`
	AppHash            string     `json:"app_hash"`
	LastResultsHash    string     `json:"last_results_hash"`
	EvidenceHash       string     `json:"evidence_hash"`
	ProposerAddress    string     `json:"proposer_address"`
}

func ToDetailsView(rawBlock *blockpb.Block) (*DetailsView, error) {
	return &DetailsView{
		AppVersion:         rawBlock.GetHeader().GetVersion().GetApp(),
		BlockVersion:       rawBlock.GetHeader().GetVersion().GetBlock(),
		ChainId:            rawBlock.GetHeader().GetChainId(),
		Height:             rawBlock.GetHeader().GetHeight(),
		Time:               *types.NewTimeFromTimestamp(*rawBlock.GetHeader().GetTime()),
		LastBlockIdHash:    rawBlock.GetHeader().GetLastBlockId().GetHash(),
		LastCommitHash:     rawBlock.GetHeader().GetLastCommitHash(),
		DataHash:           rawBlock.GetHeader().GetDataHash(),
		ValidatorsHash:     rawBlock.GetHeader().GetValidatorsHash(),
		NextValidatorsHash: rawBlock.GetHeader().GetNextValidatorsHash(),
		ConsensusHash:      rawBlock.GetHeader().GetConsensusHash(),
		AppHash:            rawBlock.GetHeader().GetAppHash(),
		LastResultsHash:    rawBlock.GetHeader().GetLastResultsHash(),
		EvidenceHash:       rawBlock.GetHeader().GetEvidenceHash(),
		ProposerAddress:    rawBlock.GetHeader().GetProposerAddress(),
	}, nil
}
