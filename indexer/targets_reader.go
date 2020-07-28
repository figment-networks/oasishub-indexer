package indexer

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/figment-networks/indexing-engine/pipeline"
	"io/ioutil"
)

const (
	TargetIndexBlockSequences = iota + 1
	TargetIndexValidatorSequences
	TargetIndexValidatorAggregates
)

var (
	_ TargetsReader = (*targetsReader)(nil)
)

type TargetsReader interface {
	GetCurrentVersionID() int64
	GetAllVersionedVersionIds() []int64
	GetAllAvailableTasks() []pipeline.TaskName
	GetAllVersionedTasks() ([]pipeline.TaskName, error)
	GetTasksByVersionIds([]int64) ([]pipeline.TaskName, error)
	GetTasksByTargetIds([]int64) ([]pipeline.TaskName, error)
}

// NewTargetsReader constructor for targetsReader
func NewTargetsReader(file string) (*targetsReader, error) {
	p := &targetsReader{}

	cfg, err := p.parseFile(file)
	if err != nil {
		return nil, err
	}

	p.cfg = cfg

	return p, nil
}

// targetsReader
type targetsReader struct {
	cfg *targetsCfg
}

type targetsCfg struct {
	Version          int64               `json:"version"`
	Versions         []version           `json:"versions"`
	SharedTasks      []pipeline.TaskName `json:"shared_tasks"`
	AvailableTargets []target            `json:"available_targets"`
}

type version struct {
	ID      int64   `json:"id"`
	Targets []int64 `json:"targets"`
}

type target struct {
	ID    int64               `json:"id"`
	Name  string              `json:"name"`
	Desc  string              `json:"desc"`
	Tasks []pipeline.TaskName `json:"tasks"`
}

//GetCurrentVersionID gets the most recent version id
func (p *targetsReader) GetCurrentVersionID() int64 {
	lastVersion := p.cfg.Versions[len(p.cfg.Versions)-1]
	return lastVersion.ID
}

// GetAllAvailableTasks get lists of tasks for all available targets
func (p *targetsReader) GetAllAvailableTasks() []pipeline.TaskName {
	var allAvailableTaskNames []pipeline.TaskName

	allAvailableTaskNames = p.appendSharedTasks(allAvailableTaskNames)

	for _, t := range p.cfg.AvailableTargets {
		allAvailableTaskNames = append(allAvailableTaskNames, t.Tasks...)
	}

	return getUniqueTaskNames(allAvailableTaskNames)
}

// GetAllVersionedVersionIds gets a slice with all version ids in the targets file
func (p *targetsReader) GetAllVersionedVersionIds() []int64 {
	currentVersionId := p.GetCurrentVersionID()
	var ids []int64
	for i := int64(1); i <= currentVersionId; i++ {
		ids = append(ids, i)
	}
	return ids
}

// GetAllVersionedTasks get lists of tasks for provided versions
func (p *targetsReader) GetAllVersionedTasks() ([]pipeline.TaskName, error) {
	var allAvailableTaskNames []pipeline.TaskName

	allAvailableTaskNames = p.appendSharedTasks(allAvailableTaskNames)

	ids := p.GetAllVersionedVersionIds()

	versionedTaskNames, err := p.GetTasksByVersionIds(ids)
	if err != nil {
		return nil, err
	}

	allAvailableTaskNames = append(allAvailableTaskNames, versionedTaskNames...)

	return getUniqueTaskNames(allAvailableTaskNames), nil
}

// GetTasksByTargetIds get lists of tasks for specific version ids
func (p *targetsReader) GetTasksByVersionIds(versionIds []int64) ([]pipeline.TaskName, error) {
	var allTaskNames []pipeline.TaskName

	allTaskNames = p.appendSharedTasks(allTaskNames)

	for _, t := range versionIds {
		tasks, err := p.getTasksByVersionId(t)
		if err != nil {
			return nil, err
		}
		allTaskNames = append(allTaskNames, tasks...)
	}

	return getUniqueTaskNames(allTaskNames), nil
}

// getTasksByVersionId get lists of tasks for specific version id
func (p *targetsReader) getTasksByVersionId(versionID int64) ([]pipeline.TaskName, error) {
	var targetIds []int64
	versionFound := false
	for _, version := range p.cfg.Versions {
		if version.ID == versionID {
			targetIds = version.Targets
			versionFound = true
		}
	}

	if !versionFound {
		return nil, errors.New(fmt.Sprintf("version %d not found", versionID))
	}

	return p.GetTasksByTargetIds(targetIds)
}

// GetTasksByTargetIds get lists of tasks for specific target ids
func (p *targetsReader) GetTasksByTargetIds(targetIds []int64) ([]pipeline.TaskName, error) {
	var allTaskNames []pipeline.TaskName

	allTaskNames = p.appendSharedTasks(allTaskNames)

	for _, t := range targetIds {
		tasks, err := p.getTasksByTargetId(t)
		if err != nil {
			return nil, err
		}
		allTaskNames = append(allTaskNames, tasks...)
	}

	return getUniqueTaskNames(allTaskNames), nil
}

// getTasksByTargetId get list of tasks for desired target id
func (p *targetsReader) getTasksByTargetId(targetId int64) ([]pipeline.TaskName, error) {
	for _, t := range p.cfg.AvailableTargets {
		if t.ID == targetId {
			return getUniqueTaskNames(t.Tasks), nil
		}
	}
	return nil, errors.New(fmt.Sprintf("target id %d does not exists", targetId))
}

// parseFile gets tasks from json files from given directory
func (p *targetsReader) parseFile(file string) (*targetsCfg, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var cfg targetsCfg
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

// appendSharedTasks appends shared tasks
func (p *targetsReader) appendSharedTasks(tasks []pipeline.TaskName) []pipeline.TaskName {
	tasks = append(tasks, p.cfg.SharedTasks...)
	return tasks
}

// getUniqueTaskNames return slice with unique task names
func getUniqueTaskNames(slice []pipeline.TaskName) []pipeline.TaskName {
	keys := make(map[pipeline.TaskName]bool)
	var list []pipeline.TaskName
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

