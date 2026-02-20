package openvswitch

import (
	"testing"
)

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	if Provider() == nil {
		t.Fatal("provider is nil")
	}
}
