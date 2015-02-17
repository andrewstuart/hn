package shell

import (
	"bytes"
	"testing"
)

func TestCli(t *testing.T) {
	b := &bytes.Buffer{}

	cli, err := NewCli(b)

	if err != nil {
		t.Errorf("Error making cli: %v", err)
	}

	cli.SetContent("foo")

	if b.String() != "foo" {
		t.Errorf("Wrong string for SetContent")
	}
}
