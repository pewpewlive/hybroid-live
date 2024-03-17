package lexer

import (
	"os"
	"testing"
)

func readFile(fileName string, t *testing.T) []byte {
	file, err := os.ReadFile(fileName)
	if err != nil {
		t.Errorf("Test failed reading file: %s\n", err.Error())
	}

	return file
}

func TestBasic(t *testing.T) {
	lexer := New(readFile("./tests/test_basic.hyb", t))
	lexer.Tokenize()
	for _, err := range lexer.Errors {
		t.Errorf("Test failed tokenizing: %v\n", err)
	}

	tokenMatch := []Token{
		{At, "@", "", 1},
		{Identifier, "Environment", "", 1},
		{LeftParen, "(", "", 1},
		{Identifier, "Level", "", 1},
		{RightParen, ")", "", 1},
		{Eof, "", "", 3},
	}

	for i, token := range lexer.Tokens {
		if token != tokenMatch[i] {
			t.Errorf("Test case failed: expected %v, got %v", tokenMatch[i].ToString(), token.ToString())
			t.FailNow()
		}
	}
}

func TestTickBlocks(t *testing.T) {
	lexer := New(readFile("./tests/test_tick_blocks.hyb", t))

	lexer.Tokenize()
	for _, err := range lexer.Errors {
		t.Errorf("Test failed tokenizing: %v\n", err)
	}

	tokenMatch := []Token{
		{At, "@", "", 1},
		{Identifier, "Environment", "", 1},
		{LeftParen, "(", "", 1},
		{Identifier, "Level", "", 1},
		{RightParen, ")", "", 1},
		{Tick, "tick", "", 3},
		{LeftBrace, "{", "", 3},
		{RightBrace, "}", "", 3},
		{Tick, "tick", "", 5},
		{With, "with", "", 5},
		{Identifier, "time", "", 5},
		{LeftBrace, "{", "", 5},
		{RightBrace, "}", "", 5},
		{Eof, "", "", 7},
	}

	for i, token := range lexer.Tokens {
		if token != tokenMatch[i] {
			t.Errorf("Test case failed: expected %v, got %v", tokenMatch[i].ToString(), token.ToString())
			t.FailNow()
		}
	}
}

func TestMatchStatement(t *testing.T) {
	lexer := New(readFile("./tests/test_match_statement.hyb", t))
	lexer.Tokenize()
	for _, err := range lexer.Errors {
		t.Errorf("Test failed tokenizing: %v\n", err)
	}

	tokenMatch := []Token{
		{At, "@", "", 1},
		{Identifier, "Environment", "", 1},
		{LeftParen, "(", "", 1},
		{Identifier, "Generic", "", 1},
		{RightParen, ")", "", 1},
		{Let, "let", "", 3},
		{Identifier, "a", "", 3},
		{Equal, "=", "", 3},
		{String, "\"something\"", "something", 3},
		{Match, "match", "", 5},
		{Identifier, "a", "", 5},
		{LeftBrace, "{", "", 5},
		{String, "\"abc\"", "abc", 6},
		{FatArrow, "=>", "", 6},
		{LeftBrace, "{", "", 6},
		{Identifier, "Print", "", 7},
		{LeftParen, "(", "", 7},
		{String, "\"Hmm\"", "Hmm", 7},
		{RightParen, ")", "", 7},
		{RightBrace, "}", "", 8},
		{String, "\"something\"", "something", 9},
		{FatArrow, "=>", "", 9},
		{LeftBrace, "{", "", 9},
		{Identifier, "Print", "", 10},
		{LeftParen, "(", "", 10},
		{String, "\"Something!\"", "Something!", 10},
		{RightParen, ")", "", 10},
		{RightBrace, "}", "", 11},
		{Identifier, "_", "", 12},
		{FatArrow, "=>", "", 12},
		{LeftBrace, "{", "", 12},
		{Identifier, "Print", "", 13},
		{LeftParen, "(", "", 13},
		{String, "\"Other!?\"", "Other!?", 13},
		{RightParen, ")", "", 13},
		{RightBrace, "}", "", 14},
		{RightBrace, "}", "", 15},
		{Eof, "", "", 17},
	}

	for i, token := range lexer.Tokens {
		if token != tokenMatch[i] {
			t.Errorf("Test case failed: expected %v, got %v", tokenMatch[i].ToString(), token.ToString())
			t.FailNow()
		}
	}
}

func TestNumberLiterals(t *testing.T) {
	lexer := New(readFile("./tests/test_number_literals.hyb", t))
	lexer.Tokenize()
	for _, err := range lexer.Errors {
		t.Errorf("Test failed tokenizing: %v\n", err)
	}

	tokenMatch := []Token{
		{Let, "let", "", 1},
		{Identifier, "a", "", 1},
		{Equal, "=", "", 1},
		{Number, "100.5f", "100.5", 1},
		{Let, "let", "", 2},
		{Identifier, "b", "", 2},
		{Equal, "=", "", 2},
		{Degree, "180d", "180", 2},
		{Let, "let", "", 3},
		{Identifier, "c", "", 3},
		{Equal, "=", "", 3},
		{Radian, "3.14r", "3.14", 3},
		{Let, "let", "", 4},
		{Identifier, "d", "", 4},
		{Equal, "=", "", 4},
		{FixedPoint, "50fx", "50", 4},
		{Let, "let", "", 5},
		{Identifier, "e", "", 5},
		{Equal, "=", "", 5},
		{Number, "6.28", "6.28", 5},
		{Let, "let", "", 6},
		{Identifier, "f", "", 6},
		{Equal, "=", "", 6},
		{Number, "5", "5", 6},
		{Let, "let", "", 7},
		{Identifier, "g", "", 7},
		{Equal, "=", "", 7},
		{Degree, "90.5d", "90.5", 7},
		{Let, "let", "", 8},
		{Identifier, "h", "", 8},
		{Equal, "=", "", 8},
		{Radian, "1r", "1", 8},
		{Let, "let", "", 9},
		{Identifier, "i", "", 9},
		{Equal, "=", "", 9},
		{Number, "42f", "42", 9},
		{Let, "let", "", 10},
		{Identifier, "j", "", 10},
		{Equal, "=", "", 10},
		{Number, "0xff00ff", "0xff00ff", 10},
		{Let, "let", "", 11},
		{Identifier, "k", "", 11},
		{Equal, "=", "", 11},
		{Number, "0xf", "0xf", 11},
		{Let, "let", "", 12},
		{Identifier, "l", "", 12},
		{Equal, "=", "", 12},
		{Number, "0xff0000ff", "0xff0000ff", 12},
		{Eof, "", "", 14},
	}

	for i, token := range lexer.Tokens {
		if token != tokenMatch[i] {
			t.Errorf("Test case failed: expected %v, got %v", tokenMatch[i].ToString(), token.ToString())
			t.FailNow()
		}
	}
}

func TestStringLiterals(t *testing.T) {
	lexer := New(readFile("./tests/test_string_literals.hyb", t))
	lexer.Tokenize()
	for _, err := range lexer.Errors {
		t.Errorf("Test failed tokenizing: %v\n", err)
	}

	tokenMatch := []Token{
		{String, "\"string literals!!!!!\"", "string literals!!!!!", 1},
		{String, "\"very cool\"", "very cool", 2},
		{String, "\"/* yes and */\"", "/* yes and */", 3},
		{String, "\"// something trait\"", "// something trait", 4},
		{String, "\"+ ^ * in to match - == != 10 ==\"", "+ ^ * in to match - == != 10 ==", 5},
		{String, "\"\\\\\\\"\\\"\\\"\\\"\"", "\\\\\\\"\\\"\\\"\\\"", 6},
		{Eof, "", "", 8},
	}

	for i, token := range lexer.Tokens {
		if token != tokenMatch[i] {
			t.Errorf("Test case failed: expected %v, got %v", tokenMatch[i].ToString(), token.ToString())
			t.FailNow()
		}
	}
}
