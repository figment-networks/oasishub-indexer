package indexer

import (
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/figment-networks/oasishub-indexer/utils/logger"
)

func TestMain(m *testing.M) {
	setup()
	exitVal := m.Run()
	os.Exit(exitVal)
}

func setup() {
	rand.Seed(time.Now().UnixNano())
	logger.InitTest()
}
