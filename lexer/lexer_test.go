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
	tokens := Tokenize(readFile("./tests/test_basic.lc", t))
	// if err != nil {
	// 	t.Errorf("Test failed tokenizing: %s\n", err.Error())
	// }

	tokenMatch := []Token{
		Token{At, "@", "", 1},
		Token{Identifier, "Environment", "", 1},
		Token{LeftParen, "(", "", 1},
		Token{Identifier, "Level", "", 1},
		Token{RightParen, ")", "", 1},
	}

	for i, token := range tokens {
		if token != tokenMatch[i] {
			t.Errorf("Test case failed: expected %v, got %v", tokenMatch[i], token)
		}
	}
}

func TestTickBlocks(t *testing.T) {
	tokens := Tokenize(readFile("./tests/test_tick_blocks.lc", t))
	// if err != nil {
	// 	t.Errorf("Test failed tokenizing: %s\n", err.Error())
	// }

	tokenMatch := []Token{
		Token{At, "@", "", 1},
		Token{Identifier, "Environment", "", 1},
		Token{LeftParen, "(", "", 1},
		Token{Identifier, "Level", "", 1},
		Token{RightParen, ")", "", 1},
		Token{Tick, "tick", "", 3},
		Token{LeftBrace, "{", "", 3},
		Token{RightBrace, "}", "", 3},
		Token{Tick, "tick", "", 5},
		Token{With, "with", "", 5},
		Token{Identifier, "time", "", 5},
		Token{LeftBrace, "{", "", 5},
		Token{RightBrace, "}", "", 5},
	}

	for i, token := range tokens {
		if token != tokenMatch[i] {
			t.Errorf("Test case failed: expected %v, got %v", tokenMatch[i], token)
		}
	}
}

func TestMatchStatement(t *testing.T) {
	tokens := Tokenize(readFile("./tests/test_match_statement.lc", t))
	// if err != nil {
	// 	t.Errorf("Test failed tokenizing: %s\n", err.Error())
	// }

	tokenMatch := []Token{
		Token{At, "@", "", 1},
		Token{Identifier, "Environment", "", 1},
		Token{LeftParen, "(", "", 1},
		Token{Identifier, "Level", "", 1},
		Token{RightParen, ")", "", 1},
		Token{Let, "let", "", 3},
		Token{Identifier, "a", "", 3},
		Token{Equal, "=", "", 3},
		Token{String, "\"something\"", "something", 3},
		Token{Match, "match", "", 5},
		Token{Identifier, "a", "", 5},
		Token{LeftBrace, "{", "", 5},
		Token{String, "\"abc\"", "abc", 6},
		Token{FatArrow, "=>", "", 6},
		Token{LeftBrace, "{", "", 6},
		Token{Identifier, "Print", "", 7},
		Token{LeftParen, "(", "", 7},
		Token{String, "\"Hmm\"", "Hmm", 7},
		Token{RightParen, ")", "", 7},
		Token{RightBrace, "}", "", 8},
		Token{String, "\"something\"", "something", 9},
		Token{FatArrow, "=>", "", 9},
		Token{LeftBrace, "{", "", 9},
		Token{Identifier, "Print", "", 10},
		Token{LeftParen, "(", "", 10},
		Token{String, "\"Something!\"", "Something!", 10},
		Token{RightParen, ")", "", 10},
		Token{RightBrace, "}", "", 11},
		Token{Identifier, "_", "", 12},
		Token{FatArrow, "=>", "", 12},
		Token{LeftBrace, "{", "", 12},
		Token{Identifier, "Print", "", 13},
		Token{LeftParen, "(", "", 13},
		Token{String, "\"Other!?\"", "Other!?", 13},
		Token{RightParen, ")", "", 13},
		Token{RightBrace, "}", "", 14},
		Token{RightBrace, "}", "", 15},
	}

	for i, token := range tokens {
		if token != tokenMatch[i] {
			t.Errorf("Test case failed: expected %v, got %v", tokenMatch[i], token)
		}
	}
}

func TestNumberLiterals(t *testing.T) {
	tokens := Tokenize(readFile("./tests/test_number_literals.lc", t))
	// if err != nil {
	// 	t.Errorf("Test failed tokenizing: %s\n", err.Error())
	// }

	tokenMatch := []Token{
		Token{Let, "let", "", 1},
		Token{Identifier, "a", "", 1},
		Token{Equal, "=", "", 1},
		Token{Number, "100.5f", "100.5", 1},
		Token{Let, "let", "", 2},
		Token{Identifier, "b", "", 2},
		Token{Equal, "=", "", 2},
		Token{Degree, "180d", "180", 2},
		Token{Let, "let", "", 3},
		Token{Identifier, "c", "", 3},
		Token{Equal, "=", "", 3},
		Token{Radian, "3.14r", "3.14", 3},
		Token{Let, "let", "", 4},
		Token{Identifier, "d", "", 4},
		Token{Equal, "=", "", 4},
		Token{FixedPoint, "50fx", "50", 4},
		Token{Let, "let", "", 5},
		Token{Identifier, "e", "", 5},
		Token{Equal, "=", "", 5},
		Token{Number, "6.28", "6.28", 5},
		Token{Let, "let", "", 6},
		Token{Identifier, "f", "", 6},
		Token{Equal, "=", "", 6},
		Token{Number, "5", "5", 6},
		Token{Let, "let", "", 7},
		Token{Identifier, "g", "", 7},
		Token{Equal, "=", "", 7},
		Token{Degree, "90.5d", "90.5", 7},
		Token{Let, "let", "", 8},
		Token{Identifier, "h", "", 8},
		Token{Equal, "=", "", 8},
		Token{Radian, "1r", "1", 8},
		Token{Let, "let", "", 9},
		Token{Identifier, "i", "", 9},
		Token{Equal, "=", "", 9},
		Token{Number, "42f", "42", 9},
		Token{Let, "let", "", 10},
		Token{Identifier, "j", "", 10},
		Token{Equal, "=", "", 10},
		Token{Number, "0xff00ff", "0xff00ff", 10},
		Token{Let, "let", "", 11},
		Token{Identifier, "k", "", 11},
		Token{Equal, "=", "", 11},
		Token{Number, "0xf", "0xf", 11},
		Token{Let, "let", "", 12},
		Token{Identifier, "l", "", 12},
		Token{Equal, "=", "", 12},
		Token{Number, "0xff0000ff", "0xff0000ff", 12},
	}

	for i, token := range tokens {
		if token != tokenMatch[i] {
			t.Errorf("Test case failed: expected %v, got %v", tokenMatch[i], token)
		}
	}
}

func TestStringLiterals(t *testing.T) {
	tokens := Tokenize(readFile("./tests/test_basic.lc", t))
	// if err != nil {
	// 	t.Errorf("Test failed tokenizing: %s\n", err.Error())
	// }

	tokenMatch := []Token{
		Token{Let, "let", "", 1},
		Token{Identifier, "a", "", 1},
		Token{Equal, "=", "", 1},
		Token{String, "\"string literals!!!!!\"", "string literals!!!!!", 1},
		Token{Let, "let", "", 2},
		Token{Identifier, "b", "", 2},
		Token{Equal, "=", "", 2},
		Token{String, "\"very cool\"", "very cool", 2},
		Token{Let, "let", "", 3},
		Token{Identifier, "c", "", 3},
		Token{Equal, "=", "", 3},
		Token{String, "\"/* yes and */\"", "/* yes and */", 3},
		Token{Let, "let", "", 4},
		Token{Identifier, "d", "", 4},
		Token{Equal, "=", "", 4},
		Token{String, "\"// something trait\"", "// something trait", 4},
		Token{Let, "let", "", 5},
		Token{Identifier, "e", "", 5},
		Token{Equal, "=", "", 5},
		Token{String, "\"+ ^ * in to match - == != 10 ==\"", "+ ^ * in to match - == != 10 ==", 5},
	}

	for i, token := range tokens {
		if token != tokenMatch[i] {
			t.Errorf("Test case failed: expected %v, got %v", tokenMatch[i], token)
		}
	}
}
