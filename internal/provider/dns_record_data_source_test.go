package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDNSRecordDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create a record first, then read it via data source
			{
				Config: testAccDNSRecordDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Check the resource
					resource.TestCheckResourceAttr("porkbun_dns_record.test_ds", "domain", testDomain),
					resource.TestCheckResourceAttr("porkbun_dns_record.test_ds", "name", "tftest-datasource"),
					resource.TestCheckResourceAttr("porkbun_dns_record.test_ds", "type", "A"),
					resource.TestCheckResourceAttr("porkbun_dns_record.test_ds", "content", "192.0.2.50"),
					// Check the data source matches the resource
					resource.TestCheckResourceAttrPair(
						"data.porkbun_dns_record.test_ds", "id",
						"porkbun_dns_record.test_ds", "id",
					),
					resource.TestCheckResourceAttrPair(
						"data.porkbun_dns_record.test_ds", "type",
						"porkbun_dns_record.test_ds", "type",
					),
					resource.TestCheckResourceAttrPair(
						"data.porkbun_dns_record.test_ds", "content",
						"porkbun_dns_record.test_ds", "content",
					),
					resource.TestCheckResourceAttrPair(
						"data.porkbun_dns_record.test_ds", "ttl",
						"porkbun_dns_record.test_ds", "ttl",
					),
				),
			},
		},
	})
}

func testAccDNSRecordDataSourceConfig() string {
	return providerConfig + fmt.Sprintf(`
resource "porkbun_dns_record" "test_ds" {
  domain  = %[1]q
  name    = "tftest-datasource"
  type    = "A"
  content = "192.0.2.50"
  ttl     = "600"
}

data "porkbun_dns_record" "test_ds" {
  domain = %[1]q
  id     = porkbun_dns_record.test_ds.id
}
`, testDomain)
}
