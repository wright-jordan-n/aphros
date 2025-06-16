package lexer

import "fmt"

type TokenTag int

const (
	TOK_EOF = iota

	TOK_INT_LIT
	TOK_HEX_LIT
	TOK_OCT_LIT
	TOK_BIN_LIT
	TOK_FLOAT_LIT

	TOK_CHAR_LIT
	TOK_STR_LIT

	TOK_IDENT

	TOK_KW_AS
	TOK_KW_ANY
	TOK_KW_BOOL
	TOK_KW_BREAK
	TOK_KW_CONTINUE
	TOK_KW_ELIF
	TOK_KW_ELSE
	TOK_KW_EXPORT
	TOK_KW_FALSE
	TOK_KW_FLOAT
	TOK_KW_FN
	TOK_KW_IF
	TOK_KW_IMPORT
	TOK_KW_INT
	TOK_KW_LOOP
	TOK_KW_MODULE
	TOK_KW_NIL
	TOK_KW_RET
	TOK_KW_STR
	TOK_KW_STRUCT
	TOK_KW_TRUE
	TOK_KW_VAR

	TOK_RSHIFT_LOGICAL_EQ
	TOK_RSHIFT_LOGICAL
	TOK_RSHIFT_ARITH_EQ
	TOK_RSHIFT_ARITH
	TOK_GT_EQ
	TOK_GT

	TOK_LSHIFT_EQ
	TOK_LSHIFT
	TOK_LT_EQ
	TOK_LT

	TOK_STAR_EQ
	TOK_STAR

	TOK_SLASH_EQ
	TOK_SLASH

	TOK_MOD_EQ
	TOK_MOD

	TOK_PLUS_EQ
	TOK_PLUS

	TOK_MINUS_EQ
	TOK_MINUS

	TOK_EQ_EQ
	TOK_EQ

	TOK_NOT_EQ
	TOK_NOT

	TOK_AND_EQ
	TOK_AND

	TOK_XOR_EQ
	TOK_XOR

	TOK_PIPE_EQ
	TOK_PIPE

	TOK_DOT

	TOK_LBRACE
	TOK_RBRACE

	TOK_LPAREN
	TOK_RPAREN

	TOK_LBRACK
	TOK_RBRACK

	TOK_COMMA
	TOK_SEMI

	TOK_ERROR
)

type Token struct {
	Tag     TokenTag
	Col     int
	Line    int
	Literal []byte
}

type LexerState struct {
	col       int
	line      int
	ok        bool
	buf       []byte
	i         int
	start     int
	col_start int
}

func advance(state *LexerState) byte {
	if state.i+1 >= len(state.buf) {
		return 0
	}
	state.i += 1
	state.col += 1
	return state.buf[state.i]
}

func newTok(state *LexerState, tag TokenTag) *Token {
	return &Token{tag, state.col_start, state.line, state.buf[state.start:state.i]}
}

func errorTok(state *LexerState, msg string) *Token {
	return &Token{TOK_ERROR, state.col_start, state.line, []byte(msg)}
}

var keywords map[string]TokenTag = map[string]TokenTag{
	"as":       TOK_KW_AS,
	"any":      TOK_KW_ANY,
	"bool":     TOK_KW_BOOL,
	"break":    TOK_KW_BREAK,
	"continue": TOK_KW_CONTINUE,
	"elif":     TOK_KW_ELIF,
	"else":     TOK_KW_ELSE,
	"export":   TOK_KW_EXPORT,
	"false":    TOK_KW_FALSE,
	"float":    TOK_KW_FLOAT,
	"fn":       TOK_KW_FN,
	"if":       TOK_KW_IF,
	"import":   TOK_KW_IMPORT,
	"int":      TOK_KW_INT,
	"loop":     TOK_KW_LOOP,
	"module":   TOK_KW_MODULE,
	"nil":      TOK_KW_NIL,
	"ret":      TOK_KW_RET,
	"str":      TOK_KW_STR,
	"struct":   TOK_KW_STRUCT,
	"true":     TOK_KW_TRUE,
	"var":      TOK_KW_VAR,
}

func Lex(buf []byte) []*Token {
	tokens := []*Token{}
	state := &LexerState{col: 1, line: 1, ok: true, buf: buf, i: -1}
	c := advance(state)
	for {
		state.start = state.i
		state.col_start = state.col
		switch c {
		case 0:
			tokens = append(tokens, &Token{TOK_EOF, state.col_start, state.line, nil})
			return tokens
		case '#':
			c = advance(state)
			for c != '\n' && c != 0 {
				c = advance(state)
			}
		case ' ':
			fallthrough
		case '\r':
			fallthrough
		case '\t':
			c = advance(state)
		case '\n':
			state.line += 1
			state.col = 0
			c = advance(state)
		case '>':
			c = advance(state)
			if c == '>' {
				c = advance(state)
				if c == '>' {
					c = advance(state)
					if c == '=' {
						c = advance(state)
						tokens = append(tokens, newTok(state, TOK_RSHIFT_LOGICAL_EQ))
						break
					}
					tokens = append(tokens, newTok(state, TOK_RSHIFT_LOGICAL))
					break
				}
				if c == '=' {
					c = advance(state)
					tokens = append(tokens, newTok(state, TOK_RSHIFT_ARITH_EQ))
					break
				}
				tokens = append(tokens, newTok(state, TOK_RSHIFT_ARITH))
				break
			}
			if c == '=' {
				c = advance(state)
				tokens = append(tokens, newTok(state, TOK_GT_EQ))
				break
			}
			tokens = append(tokens, newTok(state, TOK_GT))
		case '<':
			c = advance(state)
			if c == '<' {
				c = advance(state)
				if c == '=' {
					c = advance(state)
					tokens = append(tokens, newTok(state, TOK_LSHIFT_EQ))
					break
				}
				tokens = append(tokens, newTok(state, TOK_LSHIFT))
				break
			}
			if c == '=' {
				c = advance(state)
				tokens = append(tokens, newTok(state, TOK_LT_EQ))
				break
			}
			tokens = append(tokens, newTok(state, TOK_LT))
		case '*':
			c = advance(state)
			if c == '=' {
				c = advance(state)
				tokens = append(tokens, newTok(state, TOK_STAR_EQ))
				break
			}
			tokens = append(tokens, newTok(state, TOK_STAR))
		case '/':
			c = advance(state)
			if c == '=' {
				c = advance(state)
				tokens = append(tokens, newTok(state, TOK_SLASH_EQ))
				break
			}
			tokens = append(tokens, newTok(state, TOK_SLASH))
		case '%':
			c = advance(state)
			if c == '=' {
				c = advance(state)
				tokens = append(tokens, newTok(state, TOK_MOD_EQ))
				break
			}
			tokens = append(tokens, newTok(state, TOK_MOD))
		case '+':
			c = advance(state)
			if c == '=' {
				c = advance(state)
				tokens = append(tokens, newTok(state, TOK_PLUS_EQ))
				break
			}
			tokens = append(tokens, newTok(state, TOK_PLUS))
		case '-':
			c = advance(state)
			if c == '=' {
				c = advance(state)
				tokens = append(tokens, newTok(state, TOK_MINUS_EQ))
				break
			}
			tokens = append(tokens, newTok(state, TOK_MINUS))
		case '=':
			c = advance(state)
			if c == '=' {
				c = advance(state)
				tokens = append(tokens, newTok(state, TOK_EQ_EQ))
				break
			}
			tokens = append(tokens, newTok(state, TOK_EQ))
		case '!':
			c = advance(state)
			if c == '=' {
				c = advance(state)
				tokens = append(tokens, newTok(state, TOK_NOT_EQ))
				break
			}
			tokens = append(tokens, newTok(state, TOK_NOT))
		case '&':
			c = advance(state)
			if c == '=' {
				c = advance(state)
				tokens = append(tokens, newTok(state, TOK_AND_EQ))
				break
			}
			tokens = append(tokens, newTok(state, TOK_AND))
		case '^':
			c = advance(state)
			if c == '=' {
				c = advance(state)
				tokens = append(tokens, newTok(state, TOK_XOR_EQ))
				break
			}
			tokens = append(tokens, newTok(state, TOK_XOR))
		case '|':
			c = advance(state)
			if c == '=' {
				c = advance(state)
				tokens = append(tokens, newTok(state, TOK_PIPE_EQ))
				break
			}
			tokens = append(tokens, newTok(state, TOK_PIPE))
		case '.':
			c = advance(state)
			tokens = append(tokens, newTok(state, TOK_DOT))
		case '{':
			c = advance(state)
			tokens = append(tokens, newTok(state, TOK_LBRACE))
		case '}':
			c = advance(state)
			tokens = append(tokens, newTok(state, TOK_RBRACE))
		case '(':
			c = advance(state)
			tokens = append(tokens, newTok(state, TOK_LPAREN))
		case ')':
			c = advance(state)
			tokens = append(tokens, newTok(state, TOK_RPAREN))
		case '[':
			c = advance(state)
			tokens = append(tokens, newTok(state, TOK_LBRACK))
		case ']':
			c = advance(state)
			tokens = append(tokens, newTok(state, TOK_RBRACK))
		case ',':
			c = advance(state)
			tokens = append(tokens, newTok(state, TOK_COMMA))
		case ';':
			c = advance(state)
			tokens = append(tokens, newTok(state, TOK_SEMI))
		case '\'':
			c = advance((state))
			if c == 0 {
				state.ok = false
				tokens = append(tokens, errorTok(state, "Unterminated character literal"))
				break
			}
			if c < ' ' || c > '~' {
				state.ok = false
				tokens = append(tokens, errorTok(state, "Invalid character literal"))
				break
			}
			if c == '\'' {
				state.ok = false
				tokens = append(tokens, errorTok(state, "Empty character literal"))
				break
			}
			if c == '\\' {
				c = advance(state)
				if c != 'n' && c != 't' && c != 'r' && c != '\\' && c != '\'' {
					if c == 'x' {
						c = advance(state)
						if !((c >= 'A' && c <= 'F') || (c >= '0' && c <= '9')) {
							state.ok = false
							tokens = append(tokens, errorTok(state, "Invalid hexadecimal escape sequence"))
							break
						}
						c = advance(state)
						if !((c >= 'A' && c <= 'F') || (c >= '0' && c <= '9')) {
							state.ok = false
							tokens = append(tokens, errorTok(state, "Invalid hexadecimal escape sequence"))
							break
						}
					} else {
						state.ok = false
						tokens = append(tokens, errorTok(state, "Invalid character escape sequence"))
						break
					}
				}
				c = advance(state)
				if c != '\'' {
					state.ok = false
					tokens = append(tokens, errorTok(state, "Unterminated character literal"))
					break
				}
				c = advance(state)
				tokens = append(tokens, &Token{TOK_CHAR_LIT, state.col_start, state.line, state.buf[state.start+1 : state.i-1]})
				break
			}
			c = advance(state)
			if c != '\'' {
				tokens = append(tokens, errorTok(state, "Unterminated character literal"))
				break
			}
			tokens = append(tokens, &Token{TOK_CHAR_LIT, state.col_start, state.line, state.buf[state.start+1 : state.i-1]})
		case '"':
			c = advance((state))
			ok := true
			for c != '"' {
				if c == 0 {
					ok = false
					state.ok = false
					tokens = append(tokens, errorTok(state, "Unterminated string literal"))
					break
				}
				if c < ' ' || c > '~' {
					ok = false
					state.ok = false
					tokens = append(tokens, errorTok(state, "Invalid string literal character"))
					break
				}
				c = advance(state)
				if c == '\\' {
					c = advance(state)
					if c != 'n' && c != 't' && c != 'r' && c != '\\' && c != '"' {
						if c == 'x' {
							c = advance(state)
							if !((c >= 'A' && c <= 'F') || (c >= '0' && c <= '9')) {
								ok = false
								state.ok = false
								tokens = append(tokens, errorTok(state, "Invalid hexadecimal escape sequence"))
								break
							}
							c = advance(state)
							if !((c >= 'A' && c <= 'F') || (c >= '0' && c <= '9')) {
								ok = false
								state.ok = false
								tokens = append(tokens, errorTok(state, "Invalid hexadecimal escape sequence"))
								break
							}
						} else {
							ok = false
							state.ok = false
							tokens = append(tokens, errorTok(state, "Invalid string literal escape sequence"))
							break
						}
					}
					c = advance(state)
				}
			}
			if !ok {
				break
			}
			c = advance(state)
			tokens = append(tokens, &Token{TOK_STR_LIT, state.col_start, state.line, state.buf[state.start+1 : state.i-1]})
		case '0':
			c = advance(state)
			if c == '.' {
				c = advance(state)
				if c < '0' || c > '9' {
					state.ok = false
					tokens = append(tokens, errorTok(state, "Sequence `0.` must be proceeded with an integer digit"))
					break
				}
				c = advance(state)
				for c >= '0' && c <= '9' {
					c = advance(state)
				}
				if c == 'e' {
					c = advance(state)
					if c == '+' || c == '-' {
						advance(state)
					}
					if c < '0' || c > '9' {
						state.ok = false
						tokens = append(tokens, errorTok(state, "Sequence `0.[0-9][0-9]*e(+|-)?` must be proceeded with an integer digit"))
						break
					}
					c = advance(state)
					for c >= '0' && c <= '9' {
						c = advance(state)
					}
					tokens = append(tokens, newTok(state, TOK_FLOAT_LIT))
					break
				}
				tokens = append(tokens, newTok(state, TOK_FLOAT_LIT))
				break
			}
			if c == 'x' {
				c = advance(state)
				if !((c >= 'A' && c <= 'F') || (c >= '0' && c <= '9')) {
					state.ok = false
					tokens = append(tokens, errorTok(state, "Sequence `0x` must be proceeded with a hex digit"))
					break
				}
				c = advance(state)
				for (c >= 'A' && c <= 'F') || (c >= '0' && c <= '9') {
					c = advance(state)
				}
				tokens = append(tokens, newTok(state, TOK_HEX_LIT))
				break
			}
			if c == 'b' {
				c = advance(state)
				if c != '0' && c != '1' {
					state.ok = false
					tokens = append(tokens, errorTok(state, "Sequence `0b` must be proceeded with a binary digit"))
					break
				}
				c = advance(state)
				for c == '0' || c == '1' {
					c = advance(state)
				}
				tokens = append(tokens, newTok(state, TOK_BIN_LIT))
				break
			}
			if c == 'o' {
				c = advance(state)
				if c < '0' || c > '7' {
					state.ok = false
					tokens = append(tokens, errorTok(state, "Sequence `0o` must be proceeded with an octal digit"))
					break
				}
				c = advance(state)
				for c >= '0' && c <= '7' {
					c = advance(state)
				}
				tokens = append(tokens, newTok(state, TOK_OCT_LIT))
				break
			}
			if c == 'e' {
				c = advance(state)
				if c == '+' || c == '-' {
					advance(state)
				}
				if c < '0' || c > '9' {
					state.ok = false
					tokens = append(tokens, errorTok(state, "Sequence `0e(+|-)?` must be proceeded with an integer"))
					break
				}
				c = advance(state)
				for c >= '0' && c <= '9' {
					c = advance(state)
				}
				tokens = append(tokens, newTok(state, TOK_FLOAT_LIT))
				break
			}
			tokens = append(tokens, newTok(state, TOK_INT_LIT))
		default:
			if c >= '1' && c <= '9' {
				c = advance(state)
				for c >= '0' && c <= '9' {
					c = advance((state))
				}
				if c == '.' {
					c = advance(state)
					if c < '0' || c > '9' {
						state.ok = false
						tokens = append(tokens, errorTok(state, "Sequence `[1-9][0-9]*.` must be proceeded with an integer digit"))
						break
					}
					c = advance(state)
					for c >= '0' && c <= '9' {
						c = advance(state)
					}
					if c == 'e' {
						c = advance(state)
						if c == '+' || c == '-' {
							advance(state)
						}
						if c < '0' || c > '9' {
							state.ok = false
							tokens = append(tokens, errorTok(state, "Sequence [1-9][0-9]*.[0-9][0-9]*e(+|-)? must be proceeded with an integer digit"))
							break
						}
						c = advance(state)
						for c >= '0' && c <= '9' {
							c = advance(state)
						}
						tokens = append(tokens, newTok(state, TOK_FLOAT_LIT))
						break
					}
					tokens = append(tokens, newTok(state, TOK_FLOAT_LIT))
					break
				}
				if c == 'e' {
					c = advance(state)
					if c == '+' || c == '-' {
						advance(state)
					}
					if c < '0' || c > '9' {
						state.ok = false
						tokens = append(tokens, errorTok(state, "Sequence [1-9][0-9]*e(+|-)? must be proceeded with an integer digit"))
						break
					}
					c = advance(state)
					for c >= '0' && c <= '9' {
						c = advance(state)
					}
					tokens = append(tokens, newTok(state, TOK_FLOAT_LIT))
					break
				}
				tokens = append(tokens, newTok(state, TOK_INT_LIT))
				break
			}
			if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_' {
				c = advance(state)
				for (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_' || (c >= '0' && c <= '9') {
					c = advance(state)
				}
				kw, found := keywords[string(state.buf[state.start:state.i])]
				if found {
					tokens = append(tokens, newTok(state, kw))
					break
				}
				tokens = append(tokens, newTok(state, TOK_IDENT))
				break

			}
			state.ok = false
			tokens = append(tokens, errorTok(state, fmt.Sprintf("Invalid source code character %c", c)))
			c = advance(state)
		}
	}
}
