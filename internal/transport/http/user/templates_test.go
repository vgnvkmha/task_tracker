package user

import "testing"

func TestTemplatesParse(t *testing.T) {
	if Templates() == nil {
		t.Fatal("Templates returned nil")
	}
}
