package validator

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/figment-networks/oasishub-indexer/config"
	"github.com/figment-networks/oasishub-indexer/metric"
	"github.com/figment-networks/oasishub-indexer/store"
	"github.com/pkg/errors"
)

const (
	network_entities_filename = "amber_network_entities.csv"
)

var (
	colNames = []string{"amber entities", "entity id (new format)", "logo link"}

	ErrInValidFile = errors.New("unexpected file format")
)

type parseCSVUseCase struct {
	cfg *config.Config
	db  *store.Store
}

type record struct {
	entityName string
	logoURL    string
}

// NewParseCSVUseCase parses amber_network_entities.csv and updates logo_url and entity_name
// for each validator entry in file
func NewParseCSVUseCase(cfg *config.Config, db *store.Store) *parseCSVUseCase {
	return &parseCSVUseCase{
		cfg: cfg,
		db:  db,
	}
}

func (uc *parseCSVUseCase) Execute(ctx context.Context) error {
	defer metric.LogUseCaseDuration(time.Now(), "parse csv")

	file := fmt.Sprintf("%v/%v", uc.cfg.NetworkEntitiesDir, network_entities_filename)

	records, err := uc.parseNetworkEntities(file)
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

func (uc *parseCSVUseCase) parseNetworkEntities(file string) (map[string]*record, error) {
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

func (uc *parseCSVUseCase) updateValidatorAgg(addr string, data *record) error {
	val, err := uc.db.ValidatorAgg.FindBy("recent_address", addr)
	if err == store.ErrNotFound {
		return nil
	} else if err != nil {
		return err
	}

	val.LogoURL = data.logoURL
	val.EntityName = data.entityName

	err = uc.db.ValidatorAgg.CreateOrUpdate(val)
	if err != nil {
		return err
	}

	return nil
}

func (uc *parseCSVUseCase) validateHeaders(headers []string) error {
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
