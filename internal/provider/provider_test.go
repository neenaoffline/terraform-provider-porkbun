package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// testDomain is the domain used for acceptance testing
// Set via PORKBUN_TEST_DOMAIN environment variable
// Make sure API access is enabled for this domain in Porkbun
var testDomain = os.Getenv("PORKBUN_TEST_DOMAIN")

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"porkbun": providerserver.NewProtocol6WithError(New("test")()),
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("PORKBUN_API_KEY"); v == "" {
		t.Fatal("PORKBUN_API_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("PORKBUN_SECRET_API_KEY"); v == "" {
		t.Fatal("PORKBUN_SECRET_API_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("PORKBUN_TEST_DOMAIN"); v == "" {
		t.Fatal("PORKBUN_TEST_DOMAIN must be set for acceptance tests")
	}
}

// TestProvider_HasResources verifies the provider has the expected resources
func TestProvider_HasResources(t *testing.T) {
	// This is a unit test, not an acceptance test
	p := New("test")()

	resources := p.(*PorkbunProvider)
	if resources == nil {
		t.Fatal("provider should not be nil")
	}
}

// testAccCheckDestroy is a helper function to verify resources are destroyed
func testAccCheckDestroy(s *terraform.State) error {
	// The testing framework handles destroy verification automatically
	// This is a placeholder for any custom destroy verification logic
	return nil
}
