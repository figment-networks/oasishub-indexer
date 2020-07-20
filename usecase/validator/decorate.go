package validator

import (
	"context"
	"encoding/csv"
	"io"
	"os"
	"strings"
	"time"

	"github.com/figment-networks/oasishub-indexer/config"
	"github.com/figment-networks/oasishub-indexer/metric"
	"github.com/figment-networks/oasishub-indexer/model"
	"github.com/figment-networks/oasishub-indexer/store"
	"github.com/pkg/errors"
)

var (
	colNames = []string{"amber entities", "entity id (new format)", "logo link"}

	ErrInValidFile = errors.New("unexpected file format")
	ErrMissingFile = errors.New("missing file path")
)

type decorateUseCase struct {
	cfg *config.Config
	db  decorateStore
}

type decorateStore interface {
	CreateOrUpdate(val *model.ValidatorAgg) error
	FindBy(key string, value interface{}) (*model.ValidatorAgg, error)
}

type record struct {
	entityName string
	logoURL    string
}

// NewDecorateUseCase decorate validators based on file data. It parses a csv file
// containing logos, entity names and entity addresses for a validator, then updates
// the logo_url and entity_name for each entry
func NewDecorateUseCase(cfg *config.Config, db decorateStore) *decorateUseCase {
	return &decorateUseCase{
		cfg: cfg,
		db:  db,
	}
}

func (uc *decorateUseCase) Execute(ctx context.Context, file string) error {
	defer metric.LogUseCaseDuration(time.Now(), "decorate validator")

	if file == "" {
		return ErrMissingFile
	}

	records, err := uc.parseFile(file)
	if err != nil {
		return err
	}

	for addr, record := range records {
		err = uc.updateValidatorAgg(addr, record)
		if err != nil {
			return err
		}
	}

	return nil
}

func (uc *decorateUseCase) parseFile(file string) (map[string]*record, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := csv.NewReader(f)

	headers, err := reader.Read()
	if err != nil {
		return nil, err
	}

	err = uc.validateHeaders(headers)
	if err != nil {
		return nil, err
	}

	records := map[string]*record{}
	for {
		row, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return records, err
		}

		if row[1] == "" {
			continue
		}

		records[row[1]] = &record{
			entityName: row[0],
			logoURL:    row[2],
		}
	}
}

func (uc *decorateUseCase) updateValidatorAgg(addr string, data *record) error {
	val, err := uc.db.FindBy("recent_address", addr)
	if err == store.ErrNotFound {
		return nil
	} else if err != nil {
		return err
	}

	val.LogoURL = data.logoURL
	val.EntityName = data.entityName

	err = uc.db.CreateOrUpdate(val)
	if err != nil {
		return err
	}

	return nil
}

func (uc *decorateUseCase) validateHeaders(headers []string) error {
	if len(headers) != len(colNames) {
		return ErrInValidFile
	}

	for i, name := range colNames {
		if name != strings.ToLower(headers[i]) {
			return ErrInValidFile
		}
	}
	return nil
}
