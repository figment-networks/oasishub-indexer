package block

import (
	"github.com/figment-networks/oasis-rpc-proxy/grpc/block/blockpb"
	"github.com/figment-networks/oasishub-indexer/types"
)

type DetailsView struct {
	AppVersion         uint64
	BlockVersion       uint64
	ChainId            string
	Height             int64
	Time               types.Time
	LastBlockIdHash    string
	LastCommitHash     string
	DataHash           string
	ValidatorsHash     string
	NextValidatorsHash string
	ConsensusHash      string
	AppHash            string
	LastResultsHash    string
	EvidenceHash       string
	ProposerAddress    string
}

func ToDetailsView(rawBlock *blockpb.Block) (*DetailsView, error) {
	return &DetailsView{
		AppVersion: rawBlock.GetHeader().GetVersion().GetApp(),
		BlockVersion: rawBlock.GetHeader().GetVersion().GetBlock(),
		ChainId: rawBlock.GetHeader().GetChainId(),
		Height: rawBlock.GetHeader().GetHeight(),
		Time: *types.NewTimeFromTimestamp(*rawBlock.GetHeader().GetTime()),
		LastBlockIdHash: rawBlock.GetHeader().GetLastBlockId().GetHash(),
		LastCommitHash: rawBlock.GetHeader().GetLastCommitHash(),
		DataHash: rawBlock.GetHeader().GetDataHash(),
		ValidatorsHash: rawBlock.GetHeader().GetValidatorsHash(),
		NextValidatorsHash: rawBlock.GetHeader().GetNextValidatorsHash(),
		ConsensusHash: rawBlock.GetHeader().GetConsensusHash(),
		AppHash: rawBlock.GetHeader().GetAppHash(),
		LastResultsHash: rawBlock.GetHeader().GetLastResultsHash(),
		EvidenceHash: rawBlock.GetHeader().GetEvidenceHash(),
		ProposerAddress: rawBlock.GetHeader().GetProposerAddress(),
	}, nil
}
