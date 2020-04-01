package iterators

import (
	"fmt"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/pipeline"
)

type HeightIterator struct {
	pipeline.Iterator

	start   types.Height
	end     types.Height
	current types.Height
	err     error
}

func NewHeightIterator(start types.Height, end types.Height) *HeightIterator {
	var err error
	if end < 0 {
		err = fmt.Errorf("'end' is %d, should be >= 0", end)
	}
	return &HeightIterator{
		start:   start,
		end:     end,
		current: start,
		err:     err,
	}
}

func (i *HeightIterator) Next() bool {
	if i.err != nil {
		return false
	}
	i.current += 1
	return i.current <= i.end
}

// Value returns current even number
func (i *HeightIterator) Value() types.Height {
	if i.err != nil || i.current > i.end {
		panic("value is not valid after iterator finished")
	}
	return i.current
}

// Err returns iteration error.
func (i *HeightIterator) Error() error {
	return i.err
}

func (i *HeightIterator) Close() error {
	return nil
}

func (i *HeightIterator) StartHeight() types.Height {
	return i.start
}

func (i *HeightIterator) EndHeight() types.Height {
	return i.end
}

func (i *HeightIterator) Length() int64 {
	return int64(i.end - i.start + 1)
}
