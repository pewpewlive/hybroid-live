package lexer

import (
	"fmt"
	"hybroid/tokens"
	"os"
	"testing"
)

type TestCheck struct {
	tokenType              tokens.TokenType
	lexeme                 string
	literal                string
	line                   int
	columnStart, columnEnd int
}

func (tc TestCheck) Check(t tokens.Token) bool {
	return t != tokens.Token{Type: tc.tokenType, Lexeme: tc.lexeme, Literal: tc.literal, Location: tokens.TokenLocation{LineStart: tc.line, LineEnd: tc.line, ColStart: tc.columnStart, ColEnd: tc.columnEnd}}
}

func (tc TestCheck) ToString() string {
	return tokens.Token{Type: tc.tokenType, Lexeme: tc.lexeme, Literal: tc.literal, Location: tokens.TokenLocation{LineStart: tc.line, ColStart: tc.columnStart, ColEnd: tc.columnEnd}}.ToString()
}

func readFile(fileName string, t *testing.T) []byte {
	file, err := os.ReadFile(fileName)
	if err != nil {
		t.Errorf("Test failed reading file: %s\n", err.Error())
	}

	return file
}

func printTokens(lexer Lexer) {
	for _, token := range lexer.Tokens {
		fmt.Printf("{%v, \"%v\", \"%v\", %v, %v, %v},\n", string(token.Type), token.Lexeme, token.Literal, token.Location.LineStart, token.Location.ColStart, token.Location.ColEnd)
	}
}

func TestBasic(t *testing.T) {
	lexer := NewLexer()
	lexer.AssignSource(readFile("./tests/test_basic.hyb", t))

	lexer.Tokenize()
	for _, err := range lexer.Alerts {
		t.Errorf("Test failed tokenizing: %v\n", err)
	}

	testChecks := []TestCheck{
		{tokens.At, "@", "", 1, 1, 2},
		{tokens.Identifier, "Environment", "", 1, 2, 13},
		{tokens.LeftParen, "(", "", 1, 13, 14},
		{tokens.Identifier, "Level", "", 1, 14, 19},
		{tokens.RightParen, ")", "", 1, 19, 20},
		{tokens.Eof, "", "", 3, 60, 60},
	}

	for i, token := range lexer.Tokens {
		if testChecks[i].Check(token) {
			t.Errorf("Test case failed: expected %v, got %v", testChecks[i].ToString(), token.ToString())
			printTokens(lexer)
			t.FailNow()
		}
	}
}

func TestTickBlocks(t *testing.T) {
	lexer := NewLexer()
	lexer.AssignSource(readFile("./tests/test_tick_blocks.hyb", t))

	lexer.Tokenize()
	for _, err := range lexer.Alerts {
		t.Errorf("Test failed tokenizing: %v\n", err)
	}

	testChecks := []TestCheck{
		{tokens.At, "@", "", 1, 1, 2},
		{tokens.Identifier, "Environment", "", 1, 2, 13},
		{tokens.LeftParen, "(", "", 1, 13, 14},
		{tokens.Identifier, "Level", "", 1, 14, 19},
		{tokens.RightParen, ")", "", 1, 19, 20},
		{tokens.Tick, "tick", "", 3, 1, 5},
		{tokens.LeftBrace, "{", "", 3, 6, 7},
		{tokens.RightBrace, "}", "", 3, 7, 8},
		{tokens.Tick, "tick", "", 5, 1, 5},
		{tokens.With, "with", "", 5, 6, 10},
		{tokens.Identifier, "time", "", 5, 11, 15},
		{tokens.LeftBrace, "{", "", 5, 16, 17},
		{tokens.RightBrace, "}", "", 5, 17, 18},
		{tokens.Eof, "", "", 7, 23, 23},
	}

	for i, token := range lexer.Tokens {
		if testChecks[i].Check(token) {
			t.Errorf("Test case failed: expected %v, got %v", testChecks[i].ToString(), token.ToString())
			printTokens(lexer)
			t.FailNow()
		}
	}
}

func TestMatchStatement(t *testing.T) {
	lexer := NewLexer()
	lexer.AssignSource(readFile("./tests/test_match_statement.hyb", t))

	lexer.Tokenize()
	for _, err := range lexer.Alerts {
		t.Errorf("Test failed tokenizing: %v\n", err)
	}

	testChecks := []TestCheck{
		{tokens.At, "@", "", 1, 1, 2},
		{tokens.Identifier, "Environment", "", 1, 2, 13},
		{tokens.LeftParen, "(", "", 1, 13, 14},
		{tokens.Identifier, "Generic", "", 1, 14, 21},
		{tokens.RightParen, ")", "", 1, 21, 22},
		{tokens.Let, "let", "", 3, 1, 4},
		{tokens.Identifier, "a", "", 3, 5, 6},
		{tokens.Equal, "=", "", 3, 7, 8},
		{tokens.String, "\"something\"", "something", 3, 9, 20},
		{tokens.Match, "match", "", 5, 1, 6},
		{tokens.Identifier, "a", "", 5, 7, 8},
		{tokens.LeftBrace, "{", "", 5, 9, 10},
		{tokens.String, "\"abc\"", "abc", 6, 3, 8},
		{tokens.FatArrow, "=>", "", 6, 9, 11},
		{tokens.LeftBrace, "{", "", 6, 12, 13},
		{tokens.Identifier, "Print", "", 7, 5, 10},
		{tokens.LeftParen, "(", "", 7, 10, 11},
		{tokens.String, "\"Hmm\"", "Hmm", 7, 11, 16},
		{tokens.RightParen, ")", "", 7, 16, 17},
		{tokens.RightBrace, "}", "", 8, 3, 4},
		{tokens.String, "\"something\"", "something", 9, 3, 14},
		{tokens.FatArrow, "=>", "", 9, 15, 17},
		{tokens.LeftBrace, "{", "", 9, 18, 19},
		{tokens.Identifier, "Print", "", 10, 5, 10},
		{tokens.LeftParen, "(", "", 10, 10, 11},
		{tokens.String, "\"Something!\"", "Something!", 10, 11, 23},
		{tokens.RightParen, ")", "", 10, 23, 24},
		{tokens.RightBrace, "}", "", 11, 3, 4},
		{tokens.Identifier, "_", "", 12, 3, 4},
		{tokens.FatArrow, "=>", "", 12, 5, 7},
		{tokens.LeftBrace, "{", "", 12, 8, 9},
		{tokens.Identifier, "Print", "", 13, 5, 10},
		{tokens.LeftParen, "(", "", 13, 10, 11},
		{tokens.String, "\"Other!?\"", "Other!?", 13, 11, 20},
		{tokens.RightParen, ")", "", 13, 20, 21},
		{tokens.RightBrace, "}", "", 14, 3, 4},
		{tokens.RightBrace, "}", "", 15, 1, 2},
		{tokens.Eof, "", "", 17, 23, 23},
	}

	for i, token := range lexer.Tokens {
		if testChecks[i].Check(token) {
			t.Errorf("Test case failed: expected %v, got %v", testChecks[i].ToString(), token.ToString())
			printTokens(lexer)
			t.FailNow()
		}
	}
}

func TestNumberLiterals(t *testing.T) {
	lexer := NewLexer()
	lexer.AssignSource(readFile("./tests/test_number_literals.hyb", t))

	lexer.Tokenize()
	for _, err := range lexer.Alerts {
		t.Errorf("Test failed tokenizing: %v\n", err)
	}

	testChecks := []TestCheck{
		{tokens.Let, "let", "", 1, 1, 4},
		{tokens.Identifier, "a", "", 1, 5, 6},
		{tokens.Equal, "=", "", 1, 7, 8},
		{tokens.Fixed, "100.5f", "100.5", 1, 9, 15},
		{tokens.Let, "let", "", 2, 1, 4},
		{tokens.Identifier, "b", "", 2, 5, 6},
		{tokens.Equal, "=", "", 2, 7, 8},
		{tokens.Degree, "180d", "180", 2, 9, 13},
		{tokens.Let, "let", "", 3, 1, 4},
		{tokens.Identifier, "c", "", 3, 5, 6},
		{tokens.Equal, "=", "", 3, 7, 8},
		{tokens.Radian, "3.14r", "3.14", 3, 9, 14},
		{tokens.Let, "let", "", 4, 1, 4},
		{tokens.Identifier, "d", "", 4, 5, 6},
		{tokens.Equal, "=", "", 4, 7, 8},
		{tokens.FixedPoint, "50fx", "50", 4, 9, 13},
		{tokens.Let, "let", "", 5, 1, 4},
		{tokens.Identifier, "e", "", 5, 5, 6},
		{tokens.Equal, "=", "", 5, 7, 8},
		{tokens.Number, "6.28", "6.28", 5, 9, 13},
		{tokens.Let, "let", "", 6, 1, 4},
		{tokens.Identifier, "f", "", 6, 5, 6},
		{tokens.Equal, "=", "", 6, 7, 8},
		{tokens.Number, "5", "5", 6, 9, 10},
		{tokens.Let, "let", "", 7, 1, 4},
		{tokens.Identifier, "g", "", 7, 5, 6},
		{tokens.Equal, "=", "", 7, 7, 8},
		{tokens.Degree, "90.5d", "90.5", 7, 9, 14},
		{tokens.Let, "let", "", 8, 1, 4},
		{tokens.Identifier, "h", "", 8, 5, 6},
		{tokens.Equal, "=", "", 8, 7, 8},
		{tokens.Radian, "1r", "1", 8, 9, 11},
		{tokens.Let, "let", "", 9, 1, 4},
		{tokens.Identifier, "i", "", 9, 5, 6},
		{tokens.Equal, "=", "", 9, 7, 8},
		{tokens.Fixed, "42f", "42", 9, 9, 12},
		{tokens.Let, "let", "", 10, 1, 4},
		{tokens.Identifier, "j", "", 10, 5, 6},
		{tokens.Equal, "=", "", 10, 7, 8},
		{tokens.Number, "0xff00ff", "0xff00ff", 10, 9, 17},
		{tokens.Let, "let", "", 11, 1, 4},
		{tokens.Identifier, "k", "", 11, 5, 6},
		{tokens.Equal, "=", "", 11, 7, 8},
		{tokens.Number, "0xf", "0xf", 11, 9, 12},
		{tokens.Let, "let", "", 12, 1, 4},
		{tokens.Identifier, "l", "", 12, 5, 6},
		{tokens.Equal, "=", "", 12, 7, 8},
		{tokens.Number, "0xff0000ff", "0xff0000ff", 12, 9, 19},
		{tokens.Eof, "", "", 14, 27, 27},
	}

	for i, token := range lexer.Tokens {
		if testChecks[i].Check(token) {
			t.Errorf("Test case failed: expected %v, got %v", testChecks[i].ToString(), token.ToString())
			printTokens(lexer)
			t.FailNow()
		}
	}
}

func TestStringLiterals(t *testing.T) {
	lexer := NewLexer()
	lexer.AssignSource(readFile("./tests/test_string_literals.hyb", t))

	lexer.Tokenize()
	for _, err := range lexer.Alerts {
		t.Errorf("Test failed tokenizing: %v\n", err)
	}

	testChecks := []TestCheck{
		{tokens.String, "\"string literals!!!!!\"", "string literals!!!!!", 1, 1, 23},
		{tokens.String, "\"very cool\"", "very cool", 2, 1, 12},
		{tokens.String, "\"/* yes and */\"", "/* yes and */", 3, 1, 16},
		{tokens.String, "\"// something trait\"", "// something trait", 4, 1, 21},
		{tokens.String, "\"+ ^ * in to match - == != 10 ==\"", "+ ^ * in to match - == != 10 ==", 5, 1, 34},
		{tokens.String, "\"\\\\\\\"\\\"\\\"\\\"\"", "\\\\\\\"\\\"\\\"\\\"", 6, 1, 13},
		{tokens.Eof, "", "", 8, 19, 19},
	}

	for i, token := range lexer.Tokens {
		if testChecks[i].Check(token) {
			t.Errorf("Test case failed: expected %v, got %v", testChecks[i].ToString(), token.ToString())
			printTokens(lexer)
			t.FailNow()
		}
	}
}
