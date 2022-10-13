package num

import "unicode"

// Base represents a number base.
type Base int

const (
	Auto Base = iota
	Bin
	Oct
	Dec
	Hex
	baseCount
)

// BaseID defines the string identifier of each Base.
var BaseID = [baseCount]string{"auto", "bin", "oct", "dec", "hex"}

// String returns the string identifier of the receiver.
func (b Base) String() string { return BaseID[b] }

// BaseInt defines the number base of each Base.
var BaseInt = [baseCount]int{0, 2, 8, 10, 16}

// Int returns the number base of the receiver.
func (b Base) Int() int { return BaseInt[b] }

// ToBaseSymCase converts the given rune to the case used by Base alphabets.
var ToBaseSymCase = unicode.ToLower

// BaseSym defines the alphabet of each Base.
var BaseSym = [baseCount][]rune{
	[]rune("0123456789abcdef"),
	[]rune("01"),
	[]rune("01234567"),
	[]rune("0123456789"),
	[]rune("0123456789abcdef"),
}

// Sym returns the alphabet of the receiver.
func (b Base) Sym() []rune { return BaseSym[b] }
