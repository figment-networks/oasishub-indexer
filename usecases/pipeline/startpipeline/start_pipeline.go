package startpipeline

import (
	"context"
	"fmt"
	"github.com/figment-networks/oasishub-indexer/config"
	"github.com/figment-networks/oasishub-indexer/models/report"
	"github.com/figment-networks/oasishub-indexer/repos/accountaggrepo"
	"github.com/figment-networks/oasishub-indexer/repos/blockseqrepo"
	"github.com/figment-networks/oasishub-indexer/repos/debondingdelegationseqrepo"
	"github.com/figment-networks/oasishub-indexer/repos/delegationseqrepo"
	"github.com/figment-networks/oasishub-indexer/repos/entityaggrepo"
	"github.com/figment-networks/oasishub-indexer/repos/reportrepo"
	"github.com/figment-networks/oasishub-indexer/repos/stakingseqrepo"
	"github.com/figment-networks/oasishub-indexer/repos/syncablerepo"
	"github.com/figment-networks/oasishub-indexer/repos/transactionseqrepo"
	"github.com/figment-networks/oasishub-indexer/repos/validatorseqrepo"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
	"github.com/figment-networks/oasishub-indexer/utils/iterators"
	"github.com/figment-networks/oasishub-indexer/utils/log"
)

type UseCase interface {
	Execute(context.Context, int64) errors.ApplicationError
}

type useCase struct {
	syncableDbRepo   syncablerepo.DbRepo
	syncableProxyRepo syncablerepo.ProxyRepo

	blockSeqDbRepo               blockseqrepo.DbRepo
	validatorSeqDbRepo           validatorseqrepo.DbRepo
	transactionSeqDbRepo         transactionseqrepo.DbRepo
	stakingSeqDbRepo             stakingseqrepo.DbRepo
	accountAggDbRepo             accountaggrepo.DbRepo
	delegationSeqDbRepo          delegationseqrepo.DbRepo
	debondingDelegationSeqDbRepo debondingdelegationseqrepo.DbRepo
	entityAggDbRepo              entityaggrepo.DbRepo

	reportDbRepo reportrepo.DbRepo
}

func NewUseCase(
	syncableDbRepo syncablerepo.DbRepo,
	syncableProxyRepo syncablerepo.ProxyRepo,
	blockSeqDbRepo blockseqrepo.DbRepo,
	validatorSeqDbRepo validatorseqrepo.DbRepo,
	transactionSeqDbRepo transactionseqrepo.DbRepo,
	stakingSeqDbRepo stakingseqrepo.DbRepo,
	accountAggDbRepo accountaggrepo.DbRepo,
	delegationSeqDbRepo delegationseqrepo.DbRepo,
	debondingDelegationSeqDbRepo debondingdelegationseqrepo.DbRepo,
	entityAggDbRepo entityaggrepo.DbRepo,
	reportDbRepo reportrepo.DbRepo,
) UseCase {
	return &useCase{
		syncableDbRepo:   syncableDbRepo,
		syncableProxyRepo: syncableProxyRepo,

		blockSeqDbRepo:               blockSeqDbRepo,
		validatorSeqDbRepo:           validatorSeqDbRepo,
		transactionSeqDbRepo:         transactionSeqDbRepo,
		stakingSeqDbRepo:             stakingSeqDbRepo,
		accountAggDbRepo:             accountAggDbRepo,
		delegationSeqDbRepo:          delegationSeqDbRepo,
		debondingDelegationSeqDbRepo: debondingDelegationSeqDbRepo,
		entityAggDbRepo:              entityAggDbRepo,

		reportDbRepo: reportDbRepo,
	}
}

func (uc *useCase) Execute(ctx context.Context, batchSize int64) errors.ApplicationError {
	i, err := uc.buildIterator(batchSize)
	if err != nil {
		return err
	}

	r, err := uc.createReport(i.StartHeight(), i.EndHeight())
	if err != nil {
		return err
	}

	p := NewPipeline(
		uc.syncableDbRepo,
		uc.syncableProxyRepo,
		uc.blockSeqDbRepo,
		uc.validatorSeqDbRepo,
		uc.transactionSeqDbRepo,
		uc.stakingSeqDbRepo,
		uc.accountAggDbRepo,
		uc.delegationSeqDbRepo,
		uc.debondingDelegationSeqDbRepo,
		uc.entityAggDbRepo,
		*r,
	)
	resp := p.Start(ctx, i)

	if resp.Error == nil {
		err = uc.completeReport(r, resp)
		if err != nil {
			return err
		}
	} else {
		return errors.NewErrorFromMessage(*resp.Error, errors.PipelineProcessingError)
	}

	return nil
}

/*************** Private ***************/

func (uc *useCase) buildIterator(batchSize int64) (*iterators.HeightIterator, errors.ApplicationError) {
	// Get start height.
	h, err := uc.syncableDbRepo.GetMostRecentCommonHeight()
	var startH types.Height
	if err != nil {
		if err.Status() == errors.NotFoundError {
			startH = config.FirstBlockHeight()
		} else {
			return nil, err
		}
	} else {
		startH = *h + 1
	}
	// Get end height. It is most recent height in the node
	syncableFromNode, err := uc.syncableProxyRepo.GetHead()
	if err != nil {
		return nil, err
	}
	endH := syncableFromNode.Height

	// Validation of heights
	blocksToSyncCount := int64(endH - startH)
	if blocksToSyncCount == 0 {
		log.Warn("nothing to process", log.Field("type", "blockSync"))
		return nil, errors.NewErrorFromMessage("nothing to process", errors.PipelineProcessingError)
	}

	// Make sure that batch limit is respected
	if endH.Int64()-startH.Int64() > batchSize {
		endH = types.Height(startH.Int64()+batchSize) - 1
	}

	log.Info(fmt.Sprintf("iterator: %d - %d", startH, endH))

	i := iterators.NewHeightIterator(startH, endH)

	return i, nil
}

func (uc *useCase) createReport(start types.Height, end types.Height) (*report.Model, errors.ApplicationError) {
	r := &report.Model{
		StartHeight: start,
		EndHeight:   end,
	}
	err := uc.reportDbRepo.Create(r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (uc *useCase) completeReport(report *report.Model, resp Results) errors.ApplicationError {
	report.Complete(resp.SuccessCount, resp.ErrorCount, resp.Error, resp.Details)

	return uc.reportDbRepo.Save(report)
}
