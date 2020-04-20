package fixtures

import (
	"github.com/figment-networks/oasishub-indexer/utils/projectpath"
	"io/ioutil"
	"os"
	"path"
)

func Load(filename string) []byte {
	jsonFile, err := os.Open(path.Join(projectpath.Root, "fixtures", filename))
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		panic(err)
	}

	return byteValue
}
