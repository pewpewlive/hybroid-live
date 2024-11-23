// AUTO-GENERATED, DO NOT MANUALLY MODIFY!

package alerts

import (
	"fmt"
	"hybroid/tokens"
)

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedaParenthesis struct {
	Token    tokens.Token
	Location tokens.TokenLocation
	Symbol   string
}

func (ep *ExpectedaParenthesis) GetMessage() string {
	return fmt.Sprintf("Expected '%s'", ep.Symbol)
}

func (ep *ExpectedaParenthesis) GetTokens() []tokens.Token {
	return []tokens.Token{ep.Token}
}

func (ep *ExpectedaParenthesis) GetLocations() []tokens.TokenLocation {
	return []tokens.TokenLocation{ep.Location}
}

func (ep *ExpectedaParenthesis) GetNote() string {
	return ""
}

func (ep *ExpectedaParenthesis) GetAlertType() AlertType {
	return Error
}

func (ep *ExpectedaParenthesis) GetAlertStage() AlertStage {
	return Walker
}
