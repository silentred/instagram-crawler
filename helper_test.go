package main

import (
	"reflect"
	"testing"
)

/* Test Helpers */
func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func refute(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		t.Errorf("Did not expect %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func TestExtractUserId(t *testing.T) {
	content := `":"6vyf2ekXY-","owner":{"id":"1971661116"},"caption":"I haveually lik:{"id":"1971661116"},"caption":"I love`
	id := ExtractUserId(content)
	expect(t, id, "1971661116")
}

func TestDerterminPath(t *testing.T) {
	expect(t, 1, 1)
}

func TestDownload(t *testing.T) {

}
