package openvswitch

import (
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	otfschema "github.com/opentofu/opentofu/helper/schema"
	otftf "github.com/opentofu/opentofu/terraform"
)

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func TestProviderOpenTofu(t *testing.T) {
	providerFunc := ProviderOpenTofu()
	provider := providerFunc()
	if err := provider.(*otfschema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProviderOpenTofu_impl(t *testing.T) {
	var _ otftf.ResourceProvider = ProviderOpenTofu()()
}
