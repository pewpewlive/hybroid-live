package parser

import (
	"hybroid/ast"
	"hybroid/lexer"
)

// Appends an error to the ParserErrors
func (p *Parser) error(token lexer.Token, msg string) {
	errMsg := ast.Error{Token: token, Message: msg}
	p.Errors = append(p.Errors, errMsg)
	//panic(errMsg.Message)
}

func (p *Parser) synchronize() {
	p.advance()
	for !p.isAtEnd() {
		switch p.peek().Type {
		case lexer.For, lexer.Fn, lexer.If, lexer.Repeat, lexer.Tick,
			lexer.Return, lexer.Let, lexer.While, lexer.Pub, lexer.Const,
			lexer.Break, lexer.Continue, lexer.Add, lexer.Remove:
			return
		}

		p.advance()
	}
}

func (p *Parser) isMultiComparison() bool {
	return p.match(lexer.And, lexer.Or)
}

func (p *Parser) isComparison() bool {
	return p.match(lexer.Greater, lexer.GreaterEqual, lexer.Less, lexer.LessEqual, lexer.BangEqual, lexer.EqualEqual)
}

// Creates a BinaryExpr
func (p *Parser) createBinExpr(left ast.Node, operator lexer.Token, tokenType lexer.TokenType, lexeme string, right ast.Node) ast.Node {
	valueType := p.determineValueType(left, right)
	return ast.BinaryExpr{
		Left:      left,
		Operator:  lexer.Token{Type: tokenType, Lexeme: lexeme, Literal: "", Location: operator.Location},
		Right:     right,
		ValueType: valueType,
	}
}

// Checks if the current position the parser is at is the End Of File
func (p *Parser) isAtEnd() bool {
	return p.peek().Type == lexer.Eof
}

// Advances by one into the next token and returns the previous token before advancing
func (p *Parser) advance() lexer.Token {
	t := p.tokens[p.current]
	if p.current < len(p.tokens)-1 {
		p.current++
	}
	return t
}

// Peeks into the current token or peeks at the token that is offset from the current position by the given offset
func (p *Parser) peek(offset ...int) lexer.Token {
	if offset == nil {
		return p.tokens[p.current]
	} else {
		if p.current+offset[0] >= len(p.tokens)-1 {
			return p.tokens[p.current]
		}
		return p.tokens[p.current+offset[0]]
	}
}

// Checks if the current type is the specified token type. Returns false if it's the End Of File
func (p *Parser) check(tokenType lexer.TokenType) bool {
	if p.isAtEnd() {
		return false
	}

	return p.peek().Type == tokenType
}

// Matches the given list of tokens and advances if they match.
func (p *Parser) match(types ...lexer.TokenType) bool {
	for _, tokenType := range types {
		if p.check(tokenType) {
			p.advance()
			return true
		}
	}

	return false
}

// Consumes a list of tokens, advancing if they match and returns true. Consume also advances if none of the tokens were able to match, and returns false
func (p *Parser) consume(message string, types ...lexer.TokenType) (lexer.Token, bool) {
	if p.isAtEnd() {
		token := p.peek()
		p.error(token, message)
		return token, false // error
	}
	for _, tokenType := range types {
		if p.check(tokenType) {
			return p.advance(), true
		}
	}
	token := p.advance()
	p.error(token, message)
	return token, false // error
}

// Checks if the value type is expected to be a fixedpoint
func IsFx(valueType ast.PrimitiveValueType) bool {
	return valueType == ast.FixedPoint || valueType == ast.Fixed || valueType == ast.Radian || valueType == ast.Degree
}

func (p *Parser) getOp(opEqual lexer.Token) lexer.Token {
	switch opEqual.Type {
	case lexer.PlusEqual:
		return lexer.Token{Type: lexer.Plus, Location: opEqual.Location, Literal: opEqual.Literal, Lexeme: "+"}
	case lexer.MinusEqual:
		return lexer.Token{Type: lexer.Minus, Location: opEqual.Location, Literal: opEqual.Literal, Lexeme: "-"}
	case lexer.SlashEqual:
		return lexer.Token{Type: lexer.Slash, Location: opEqual.Location, Literal: opEqual.Literal, Lexeme: "/"}
	case lexer.StarEqual:
		return lexer.Token{Type: lexer.Star, Location: opEqual.Location, Literal: opEqual.Literal, Lexeme: "*"}
	case lexer.CaretEqual:
		return lexer.Token{Type: lexer.Caret, Location: opEqual.Location, Literal: opEqual.Literal, Lexeme: "^"}
	case lexer.ModuloEqual:
		return lexer.Token{Type: lexer.Modulo, Location: opEqual.Location, Literal: opEqual.Literal, Lexeme: "%"}
	default:
		return lexer.Token{}
	}
}

func (p *Parser) getParam() ast.Param {
	paramName := p.expression()
	paramType := p.Type()
	if paramName.GetType() != ast.Identifier {
		p.error(paramName.GetToken(), "expected an identifier in parameter")
	}
	return ast.Param{Type: paramType, Name: paramName.GetToken()}
}

func (p *Parser) parameters() []ast.Param {
	if _, ok := p.consume("expected opening paren after an identifier", lexer.LeftParen); !ok {
		return nil
	}

	var args []ast.Param
	if p.match(lexer.RightParen) {
		args = make([]ast.Param, 0)
	} else {
		args = append(args, p.getParam())
		for p.match(lexer.Comma) {
			args = append(args, p.getParam())
		}
		p.consume("expected closing paren after parameters", lexer.RightParen)
	}

	return args
}

func (p *Parser) arguments() []ast.Node {
	if _, ok := p.consume("expected opening paren", lexer.LeftParen); !ok {
		return nil
	}

	var args []ast.Node
	if p.match(lexer.RightParen) {
		args = make([]ast.Node, 0)
	} else {
		arg := p.expression()
		args = append(args, arg)
		for p.match(lexer.Comma) {
			arg := p.expression()
			args = append(args, arg)
		}
		p.consume("expected closing paren after arguments", lexer.RightParen)
	}

	return args
}
