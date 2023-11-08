package main

import "fmt"

type Location struct {
	line int
	col  int
}

type Keyword string

const (
	selectK Keyword = "select"
	fromK   Keyword = "from"
	asK     Keyword = "as"
	tableK  Keyword = "table"
	createK Keyword = "create"
	insertK Keyword = "insert"
	intoK   Keyword = "into"
	valuesK Keyword = "values"
	intK    Keyword = "int"
	textK   Keyword = "text"
)

type Symbol string

const (
	semicolonS Symbol = ";"
	commaS     Symbol = ","
	starS      Symbol = "*"
	leftPareS  Symbol = "("
	rightPareS Symbol = ")"
)

type TokenKind uint

const (
	keywordT TokenKind = iota
	SymbolT
	identifierT
	stringT
	numericT
)

type Token struct {
	value string
	kind  TokenKind
	loc   Location
}

type Cursor struct {
	pointer int
	loc     Location
}

func (t *Token) equals(other *Token) bool {
	return t.value == other.value && t.kind == other.kind
}

type Lexer func(string, Cursor) (*Token, Cursor, bool)

func Lex(source string) ([]*Token, error) {
	tokens := []*Token{}
	cur := Cursor{}

Lex:
	for cur.pointer < len(source) {
		lexers := []Lexer{lexKeyword, lexString, lexSymbol, lexIdentifier, lexNumeric}

		for _, l := range lexers {
			if token, newCursor, ok := l(source, cur); ok {
				cur = newCursor

				if token != nil {
					tokens = append(tokens, token)
				}

				continue Lex
			}
		}

		hint := ""

		if len(tokens) > 0 {
			hint = " after " + tokens[len(tokens)-1].value
		}
		return nil, fmt.Errorf("Unable to lex token%s, at %d:%d", hint, cur.loc.line, cur.loc.col)

	}

	return tokens, nil
}

func lexNumeric(source string, ic Cursor) (*Token, Cursor, bool) {
	cur := ic

	periodFound := false
	expmarkerFound := false

	for ; cur.pointer < len(source); cur.pointer++ {
		c := source[cur.pointer]

		cur.loc.col++

		isDigit := c >= '0' && c <= '9'
		isPeriod := c == '.'
		isExpMarker := c == 'e' || c == 'E'

		if cur.pointer == ic.pointer {
			if !isDigit && !isPeriod {
				return nil, ic, false
			}

			periodFound = isPeriod
			continue
		}

		if isPeriod {
			if periodFound {
				return nil, ic, false
			}

			periodFound = true
			continue
		}

		if isExpMarker {
			if expmarkerFound {
				return nil, ic, false
			}

			periodFound = true
			expmarkerFound = true

			if cur.pointer == len(source)-1 {
				return nil, ic, false
			}

			cNext := source[cur.pointer+1]
			if cNext == '+' || cNext == '-' {
				cur.pointer++
				cur.loc.col++
			}

			continue
		}

		if !isDigit {
			continue
		}
	}

	if cur.pointer == ic.pointer {
		return nil, ic, false
	}

	return &Token{
		value: source[ic.pointer:cur.pointer],
		kind:  numericT,
		loc:   ic.loc,
	}, cur, true
}

func lexCharacterDelimiter(source string, ic Cursor, delimiter byte) (*Token, Cursor, bool) {
	cur := ic

	if len(source[cur.pointer:]) == 0 {
		return nil, ic, false
	}

	if source[cur.pointer] != delimiter {
		return nil, ic, false
	}

	cur.loc.col++
	cur.pointer++

	var value []byte

	for ; cur.pointer < len(source); cur.pointer++ {
		c := source[cur.pointer]

		if c == delimiter {
			if cur.pointer+1 >= len(source) || source[cur.pointer+1] != delimiter {
				return &Token{
					value: string(value),
					kind:  stringT,
					loc:   ic.loc,
				}, cur, true
			} else {
				value = append(value, delimiter)
				cur.pointer++
				cur.loc.col++
			}
		}
		value = append(value, c)
		cur.loc.col++
	}

	return nil, ic, false
}

func lexString(source string, ic Cursor) (*Token, Cursor, bool) {
	return nil, ic, false
}

func lexSymbol(source string, ic Cursor) (*Token, Cursor, bool) {
	return nil, ic, false
}

func lexKeyword(source string, ic Cursor) (*Token, Cursor, bool) {
	return nil, ic, false
}

func main() {

}
