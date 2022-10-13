package op

type Operator interface {
	Ops() int           // Number of operands
	Sym() string        // Symbolic encoding
	Enc() string        // Internal encoding
	Int() int           // Precedence level
	Cmp(p Operator) int // Precedence comparison (<0 if p has higher precedence)
}

// Parse returns the Operator that is encoded as r symbolically.
func Parse(r string, n int) (p Operator, ok bool) {
	switch n {
	case 1:
		for s, e := range unarySym {
			if r == e {
				return Unary(s), true
			}
		}
	case 2:
		for s, e := range binarySym {
			if r == e {
				return Binary(s), true
			}
		}
	}
	return
}

// Decode returns the Operator that is encoded as r internally.
func Decode(r string) (p Operator, ok bool) {
	for u, e := range unaryEnc {
		if r == e {
			return Unary(u), true
		}
	}
	for b, e := range binaryEnc {
		if r == e {
			return Binary(b), true
		}
	}
	return
}

type Binary int

const (
	Xor Binary = iota // Bitwise XOR
	Bor               // Bitwise OR
	Dif               // Difference
	Sum               // Sum
	Bcl               // Bitwise clear (AND NOT)
	Ban               // Bitwise AND
	Rsh               // Bitwise right-shift
	Lsh               // Bitwise left-shift
	Rem               // Remainder (modulo)
	Quo               // Quotient
	Pro               // Product
	Pow               // Power (exponent)
	binaryCount
)

var (
	binarySym = [binaryCount]string{
		"^", "|", "-", "+", "&^", "&", ">>", "<<", "%", "/", "*", "**",
	}
	binaryEnc = [binaryCount]string{
		"^", "|", "-", "+", "\\", "&", "}", "{", "%", "/", "*", "`",
	}
)

func (Binary) Ops() int             { return 2 }
func (b Binary) Sym() string        { return binarySym[b] }
func (b Binary) Enc() string        { return binaryEnc[b] }
func (b Binary) Int() int           { return int(b) }
func (b Binary) Cmp(p Operator) int { return b.Int() - p.Int() }

type Unary int

const (
	Neg Unary = iota // Negation
	Add              // Addition
	Cmp              // Bitwise complement
	Not              // Logical complement (NOT)
	unaryCount
)

var (
	unarySym = [unaryCount]string{"-", "+", "~", "!"}
	unaryEnc = [unaryCount]string{"_", " ", "~", "!"}
)

func (Unary) Ops() int             { return 1 }
func (u Unary) Sym() string        { return unarySym[u] }
func (u Unary) Enc() string        { return unaryEnc[u] }
func (u Unary) Int() int           { return int(u) + int(binaryCount) }
func (u Unary) Cmp(p Operator) int { return u.Int() - p.Int() }
