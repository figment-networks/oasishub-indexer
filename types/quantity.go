package types

import (
	"database/sql/driver"
	"fmt"
	"math/big"
)

var (
	zero big.Int
)

type Quantity struct {
	big.Int
}

// NewQuantity creates a new Quantity from a big.Int
func NewQuantity(i *big.Int) Quantity {
	return Quantity{Int: *i}
}

// NewQuantityFromInt64 creates a new Quantity from an int64
func NewQuantityFromInt64(i int64) Quantity {
	b := big.NewInt(i)
	return Quantity{Int: *b}
}

// NewQuantityFromBytes creates a new Quantity from bytes
func NewQuantityFromBytes(bytes []byte) Quantity {
	b := big.Int{}
	return Quantity{Int: *b.SetBytes(bytes)}
}

// Valid returns true iff b >= 0
func (b *Quantity) Valid() bool {
	return b.Int.Sign() >= 0
}

// IsZero returns true iff b equals zero
func (b *Quantity) IsZero() bool {
	return b.Int.CmpAbs(&zero) == 0
}

// Equals returns true iff b equals o
func (b *Quantity) Equals(o Quantity) bool {
	return b.Int.String() == o.Int.String()
}

// Sub subtracts exactly o from b, returning an error if b < o or o < 0
func (b *Quantity) Sub(o Quantity) error {
	if !o.Valid() {
		return fmt.Errorf("could not subtract %v: invalid quantity", o)
	}
	if b.Cmp(o) == -1 {
		return fmt.Errorf("could not subtract %v: subtrahend must be smaller than %v", o, b)
	}
	b.Int.Sub(&b.Int, &o.Int)
	return nil
}

// Mul multiplies n with q, returning an error if o < 0
func (b *Quantity) Mul(o Quantity) error {
	if !o.Valid() {
		return fmt.Errorf("could not multiply %v: invalid quantity", o)
	}
	b.Int.Mul(&b.Int, &o.Int)
	return nil
}

// Quo divides b with o, returning an error if o <= 0
func (b *Quantity) Quo(o Quantity) error {
	if !o.Valid() || o.IsZero() {
		return fmt.Errorf("could not divide %v: invalid quantity", o)
	}
	b.Int.Quo(&b.Int, &o.Int)
	return nil
}

// Cmp returns -1 if b < o, 0 if b == o, and 1 if b > o
func (b *Quantity) Cmp(o Quantity) int {
	cmpTo := &o.Int
	return b.Int.Cmp(cmpTo)
}

// Clone copies a Quantity.
func (b *Quantity) Clone() Quantity {
	tmp := Quantity{}
	tmp.Set(&b.Int)
	return tmp
}

// Value implement sql.Scanner
func (b *Quantity) Value() (driver.Value, error) {
	if b != nil {
		return (b).String(), nil
	}
	return nil, nil
}

func (b *Quantity) Scan(value interface{}) error {
	b.Int = *new(big.Int)
	if value == nil {
		return nil
	}
	switch t := value.(type) {
	case int64:
		b.SetInt64(t)
	case []byte:
		b.SetString(string(value.([]byte)), 10)
	case string:
		b.SetString(t, 10)
	default:
		return fmt.Errorf("could not scan type %T into BigInt ", t)
	}
	return nil
}
