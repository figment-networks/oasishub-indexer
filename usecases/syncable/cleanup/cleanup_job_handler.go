package cleanup

import (
	"context"
	"github.com/figment-networks/oasishub-indexer/config"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/log"
)

type jobHandler struct {
	useCase UseCase
}

func NewJobHandler(useCase UseCase) types.JobHandler {
	return &jobHandler{
		useCase: useCase,
	}
}

func (h *jobHandler) Handle() {
	threshold := config.CleanupThreshold()
	ctx := context.Background()

	err := h.useCase.Execute(ctx, threshold)
	if err != nil {
		log.Error(err)
		return
	}
}

