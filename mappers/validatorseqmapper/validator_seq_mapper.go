package validatorseqmapper

import (
	"github.com/figment-networks/oasishub/db/timescale/orm"
	"github.com/figment-networks/oasishub/domain/commons"
	"github.com/figment-networks/oasishub/domain/syncabledomain"
	"github.com/figment-networks/oasishub/domain/validatordomain"
	"github.com/figment-networks/oasishub/mappers/syncablemapper"
	"github.com/figment-networks/oasishub/types"
	"github.com/figment-networks/oasishub/utils/errors"
)

func FromPersistence(o orm.ValidatorSeqModel) (*validatordomain.ValidatorSeq, errors.ApplicationError) {
	e := &validatordomain.ValidatorSeq{
		DomainEntity: commons.NewDomainEntity(commons.EntityProps{
			ID: o.ID,
		}),
		Sequence: commons.NewSequence(commons.SequenceProps{
			ChainId: o.ChainId,
			Height:  o.Height,
			Time:    o.Time,
		}),
		EntityUID:    o.EntityUID,
		NodeUID:      o.NodeUID,
		ConsensusUID: o.ConsensusUID,
		Address:      o.Address,
		VotingPower:  o.VotingPower,
		Precommit:    &validatordomain.Precommit{},
	}

	if o.PrecommitValidated != nil {
		e.Precommit.Validated = *o.PrecommitValidated

	}
	if o.PrecommitType != nil {
		e.Precommit.Type = *o.PrecommitType

	}
	if o.PrecommitIndex != nil {
		e.Precommit.Index = *o.PrecommitIndex
	}

	if !e.Valid() {
		return nil, errors.NewErrorFromMessage("validator sequence not valid", errors.NotValid)
	}

	return e, nil
}

func ToPersistence(e *validatordomain.ValidatorSeq) (*orm.ValidatorSeqModel, errors.ApplicationError) {
	if !e.Valid() {
		return nil, errors.NewErrorFromMessage("validator sequence not valid", errors.NotValid)
	}

	precommit := e.Precommit

	v := &orm.ValidatorSeqModel{
		EntityModel: orm.EntityModel{ID: e.ID},
		SequenceModel: orm.SequenceModel{
			ChainId: e.ChainId,
			Height:  e.Height,
			Time:    e.Time,
		},

		EntityUID:    e.EntityUID,
		NodeUID:      e.NodeUID,
		ConsensusUID: e.ConsensusUID,
		Address:      e.Address,
		VotingPower:  e.VotingPower,
	}

	if precommit != nil {
		v.PrecommitValidated = &precommit.Validated
		v.PrecommitType = &precommit.Type
		v.PrecommitIndex = &precommit.Index
	}

	return v, nil
}

func FromData(validatorsSyncable syncabledomain.Syncable, blockSyncable syncabledomain.Syncable) ([]*validatordomain.ValidatorSeq, errors.ApplicationError) {
	validatorsData, err := syncablemapper.UnmarshalValidatorsData(validatorsSyncable.Data)
	if err != nil {
		return nil, err
	}
	blockData, err := syncablemapper.UnmarshalBlockData(blockSyncable.Data)
	if err != nil {
		return nil, err
	}

	var validators []*validatordomain.ValidatorSeq
	for i, rv := range validatorsData.Data {
		e := &validatordomain.ValidatorSeq{
			DomainEntity: commons.NewDomainEntity(commons.EntityProps{}),
			Sequence: commons.NewSequence(commons.SequenceProps{
				ChainId: validatorsSyncable.ChainId,
				Height:  validatorsSyncable.Height,
				Time:    validatorsSyncable.Time,
			}),

			EntityUID:    types.PublicKey(rv.Node.EntityID.String()),
			NodeUID:      types.PublicKey(rv.ID.String()),
			ConsensusUID: types.PublicKey(rv.Node.Consensus.ID.String()),
			Address:      rv.Address,
			VotingPower:  validatordomain.VotingPower(rv.VotingPower),
		}

		// Block #1 does not have precommits
		if len(blockData.Data.LastCommit.Precommits) > 0 {
			precommit := blockData.Data.LastCommit.Precommits[i]

			if precommit != nil {
				e.Precommit = &validatordomain.Precommit{}
				e.Precommit.Validated = true
				e.Precommit.Type = int64(precommit.Type)
				e.Precommit.Index = int64(precommit.ValidatorIndex)
			}
		}

		if !e.Valid() {
			return nil, errors.NewErrorFromMessage("validator sequence not valid", errors.NotValid)
		}

		validators = append(validators, e)
	}
	return validators, nil
}
