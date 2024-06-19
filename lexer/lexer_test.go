package lexer

import (
	"fmt"
	"os"
	"testing"
)

type TestCheck struct {
	tokenType              TokenType
	lexeme                 string
	literal                string
	line                   int
	columnStart, columnEnd int
}

func (tc TestCheck) Check(t Token) bool {
	return t != Token{tc.tokenType, tc.lexeme, tc.literal, TokenLocation{LineStart: tc.line, LineEnd: tc.line, ColStart: tc.columnStart, ColEnd: tc.columnEnd}}
}

func (tc TestCheck) ToString() string {
	return Token{tc.tokenType, tc.lexeme, tc.literal, TokenLocation{LineStart: tc.line, ColStart: tc.columnStart, ColEnd: tc.columnEnd}}.ToString()
}

func readFile(fileName string, t *testing.T) []byte {
	file, err := os.ReadFile(fileName)
	if err != nil {
		t.Errorf("Test failed reading file: %s\n", err.Error())
	}

	return file
}

func printTokens(lexer *Lexer) {
	for _, token := range lexer.Tokens {
		fmt.Printf("{%v, \"%v\", \"%v\", %v, %v, %v},\n", string(token.Type), token.Lexeme, token.Literal, token.Location.LineStart, token.Location.ColStart, token.Location.ColEnd)
	}
}

func TestBasic(t *testing.T) {
	lexer := NewLexer()
	lexer.AssignSource(readFile("./tests/test_basic.hyb", t))

	lexer.Tokenize()
	for _, err := range lexer.Errors {
		t.Errorf("Test failed tokenizing: %v\n", err)
	}

	testChecks := []TestCheck{
		{At, "@", "", 1, 1, 2},
		{Identifier, "Environment", "", 1, 2, 13},
		{LeftParen, "(", "", 1, 13, 14},
		{Identifier, "Level", "", 1, 14, 19},
		{RightParen, ")", "", 1, 19, 20},
		{Eof, "", "", 3, 60, 60},
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
	for _, err := range lexer.Errors {
		t.Errorf("Test failed tokenizing: %v\n", err)
	}

	testChecks := []TestCheck{
		{At, "@", "", 1, 1, 2},
		{Identifier, "Environment", "", 1, 2, 13},
		{LeftParen, "(", "", 1, 13, 14},
		{Identifier, "Level", "", 1, 14, 19},
		{RightParen, ")", "", 1, 19, 20},
		{Tick, "tick", "", 3, 1, 5},
		{LeftBrace, "{", "", 3, 6, 7},
		{RightBrace, "}", "", 3, 7, 8},
		{Tick, "tick", "", 5, 1, 5},
		{With, "with", "", 5, 6, 10},
		{Identifier, "time", "", 5, 11, 15},
		{LeftBrace, "{", "", 5, 16, 17},
		{RightBrace, "}", "", 5, 17, 18},
		{Eof, "", "", 7, 23, 23},
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
	for _, err := range lexer.Errors {
		t.Errorf("Test failed tokenizing: %v\n", err)
	}

	testChecks := []TestCheck{
		{At, "@", "", 1, 1, 2},
		{Identifier, "Environment", "", 1, 2, 13},
		{LeftParen, "(", "", 1, 13, 14},
		{Identifier, "Generic", "", 1, 14, 21},
		{RightParen, ")", "", 1, 21, 22},
		{Let, "let", "", 3, 1, 4},
		{Identifier, "a", "", 3, 5, 6},
		{Equal, "=", "", 3, 7, 8},
		{String, "\"something\"", "something", 3, 9, 20},
		{Match, "match", "", 5, 1, 6},
		{Identifier, "a", "", 5, 7, 8},
		{LeftBrace, "{", "", 5, 9, 10},
		{String, "\"abc\"", "abc", 6, 3, 8},
		{FatArrow, "=>", "", 6, 9, 11},
		{LeftBrace, "{", "", 6, 12, 13},
		{Identifier, "Print", "", 7, 5, 10},
		{LeftParen, "(", "", 7, 10, 11},
		{String, "\"Hmm\"", "Hmm", 7, 11, 16},
		{RightParen, ")", "", 7, 16, 17},
		{RightBrace, "}", "", 8, 3, 4},
		{String, "\"something\"", "something", 9, 3, 14},
		{FatArrow, "=>", "", 9, 15, 17},
		{LeftBrace, "{", "", 9, 18, 19},
		{Identifier, "Print", "", 10, 5, 10},
		{LeftParen, "(", "", 10, 10, 11},
		{String, "\"Something!\"", "Something!", 10, 11, 23},
		{RightParen, ")", "", 10, 23, 24},
		{RightBrace, "}", "", 11, 3, 4},
		{Identifier, "_", "", 12, 3, 4},
		{FatArrow, "=>", "", 12, 5, 7},
		{LeftBrace, "{", "", 12, 8, 9},
		{Identifier, "Print", "", 13, 5, 10},
		{LeftParen, "(", "", 13, 10, 11},
		{String, "\"Other!?\"", "Other!?", 13, 11, 20},
		{RightParen, ")", "", 13, 20, 21},
		{RightBrace, "}", "", 14, 3, 4},
		{RightBrace, "}", "", 15, 1, 2},
		{Eof, "", "", 17, 23, 23},
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
	for _, err := range lexer.Errors {
		t.Errorf("Test failed tokenizing: %v\n", err)
	}

	testChecks := []TestCheck{
		{Let, "let", "", 1, 1, 4},
		{Identifier, "a", "", 1, 5, 6},
		{Equal, "=", "", 1, 7, 8},
		{Fixed, "100.5f", "100.5", 1, 9, 15},
		{Let, "let", "", 2, 1, 4},
		{Identifier, "b", "", 2, 5, 6},
		{Equal, "=", "", 2, 7, 8},
		{Degree, "180d", "180", 2, 9, 13},
		{Let, "let", "", 3, 1, 4},
		{Identifier, "c", "", 3, 5, 6},
		{Equal, "=", "", 3, 7, 8},
		{Radian, "3.14r", "3.14", 3, 9, 14},
		{Let, "let", "", 4, 1, 4},
		{Identifier, "d", "", 4, 5, 6},
		{Equal, "=", "", 4, 7, 8},
		{FixedPoint, "50fx", "50", 4, 9, 13},
		{Let, "let", "", 5, 1, 4},
		{Identifier, "e", "", 5, 5, 6},
		{Equal, "=", "", 5, 7, 8},
		{Number, "6.28", "6.28", 5, 9, 13},
		{Let, "let", "", 6, 1, 4},
		{Identifier, "f", "", 6, 5, 6},
		{Equal, "=", "", 6, 7, 8},
		{Number, "5", "5", 6, 9, 10},
		{Let, "let", "", 7, 1, 4},
		{Identifier, "g", "", 7, 5, 6},
		{Equal, "=", "", 7, 7, 8},
		{Degree, "90.5d", "90.5", 7, 9, 14},
		{Let, "let", "", 8, 1, 4},
		{Identifier, "h", "", 8, 5, 6},
		{Equal, "=", "", 8, 7, 8},
		{Radian, "1r", "1", 8, 9, 11},
		{Let, "let", "", 9, 1, 4},
		{Identifier, "i", "", 9, 5, 6},
		{Equal, "=", "", 9, 7, 8},
		{Fixed, "42f", "42", 9, 9, 12},
		{Let, "let", "", 10, 1, 4},
		{Identifier, "j", "", 10, 5, 6},
		{Equal, "=", "", 10, 7, 8},
		{Number, "0xff00ff", "0xff00ff", 10, 9, 17},
		{Let, "let", "", 11, 1, 4},
		{Identifier, "k", "", 11, 5, 6},
		{Equal, "=", "", 11, 7, 8},
		{Number, "0xf", "0xf", 11, 9, 12},
		{Let, "let", "", 12, 1, 4},
		{Identifier, "l", "", 12, 5, 6},
		{Equal, "=", "", 12, 7, 8},
		{Number, "0xff0000ff", "0xff0000ff", 12, 9, 19},
		{Eof, "", "", 14, 27, 27},
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
	for _, err := range lexer.Errors {
		t.Errorf("Test failed tokenizing: %v\n", err)
	}

	testChecks := []TestCheck{
		{String, "\"string literals!!!!!\"", "string literals!!!!!", 1, 1, 23},
		{String, "\"very cool\"", "very cool", 2, 1, 12},
		{String, "\"/* yes and */\"", "/* yes and */", 3, 1, 16},
		{String, "\"// something trait\"", "// something trait", 4, 1, 21},
		{String, "\"+ ^ * in to match - == != 10 ==\"", "+ ^ * in to match - == != 10 ==", 5, 1, 34},
		{String, "\"\\\\\\\"\\\"\\\"\\\"\"", "\\\\\\\"\\\"\\\"\\\"", 6, 1, 13},
		{Eof, "", "", 8, 19, 19},
	}

	for i, token := range lexer.Tokens {
		if testChecks[i].Check(token) {
			t.Errorf("Test case failed: expected %v, got %v", testChecks[i].ToString(), token.ToString())
			printTokens(lexer)
			t.FailNow()
		}
	}
}
