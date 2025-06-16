package lexer

type TokenTag int

const (
	TOK_EOF = iota

	TOK_INT_LIT
	TOK_HEX_LIT
	TOK_OCT_LIT
	TOK_BIN_LIT
	TOK_FLOAT_LIT

	// TODO
	TOK_CHAR_LIT
	TOK_STR_LIT

	// TODO
	TOK_IDENT

	// TODO
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
	tag     TokenTag
	col     int
	line    int
	literal []byte
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

func Lex(buf []byte) []*Token {
	tokens := []*Token{}
	state := &LexerState{col: 1, line: 1, ok: true, buf: buf, i: -1}
	c := advance(state)
	for {
		state.start = state.i
		state.col_start = state.col
		switch c {
		case 0:
			tokens = append(tokens, &Token{TOK_EOF, state.col, state.line, nil})
			return tokens
		case '\n':
			state.line += 1
			state.col = 1
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
		case '0':
			c = advance(state)
			if c == '.' {
				c = advance(state)
				if c < '0' || c > '9' {
					tokens = append(tokens, newTok(state, TOK_ERROR))
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
						tokens = append(tokens, newTok(state, TOK_ERROR))
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
				if c < 'A' || c > 'F' {
					tokens = append(tokens, newTok(state, TOK_ERROR))
					break
				}
				c = advance(state)
				for c >= 'A' && c <= 'F' {
					c = advance(state)
				}
				tokens = append(tokens, newTok(state, TOK_HEX_LIT))
				break
			}
			if c == 'b' {
				c = advance(state)
				if c != '0' && c != '1' {
					tokens = append(tokens, newTok(state, TOK_ERROR))
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
					tokens = append(tokens, newTok(state, TOK_ERROR))
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
					tokens = append(tokens, newTok(state, TOK_ERROR))
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
						tokens = append(tokens, newTok(state, TOK_ERROR))
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
							tokens = append(tokens, newTok(state, TOK_ERROR))
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
					if c < 'A' || c > 'F' {
						tokens = append(tokens, newTok(state, TOK_ERROR))
						break
					}
					c = advance(state)
					for c >= 'A' && c <= 'F' {
						c = advance(state)
					}
					tokens = append(tokens, newTok(state, TOK_HEX_LIT))
					break
				}
				if c == 'b' {
					c = advance(state)
					if c != '0' && c != '1' {
						tokens = append(tokens, newTok(state, TOK_ERROR))
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
						tokens = append(tokens, newTok(state, TOK_ERROR))
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
						tokens = append(tokens, newTok(state, TOK_ERROR))
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
			c = advance(state)
			state.ok = false
			tokens = append(tokens, &Token{TOK_ERROR, state.col, state.line, state.buf[state.start:state.i]})
		}
	}
}
