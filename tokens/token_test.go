package tokens

import "testing"

func TestToken(t *testing.T) {
	if Identifier.String() != "<identifier>" {
		t.Fail()
	}
	if Token(1e6).String() != "<unknown>" {
		t.Fail()
	}
	if Left.String() != "left" {
		t.Fail()
	}
	if PseudoToken(1e6).String() != "<unknown pseudo token>" {
		t.Fail()
	}
}
