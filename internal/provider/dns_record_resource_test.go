package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// importStateIdFunc returns an ImportStateIdFunc for a given resource name
func importStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", resourceName)
		}
		return fmt.Sprintf("%s/%s", testDomain, rs.Primary.ID), nil
	}
}

func TestAccDNSRecordResource_A(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccDNSRecordResourceConfig_A("tftest", "192.0.2.1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("porkbun_dns_record.test", "domain", testDomain),
					resource.TestCheckResourceAttr("porkbun_dns_record.test", "name", "tftest"),
					resource.TestCheckResourceAttr("porkbun_dns_record.test", "type", "A"),
					resource.TestCheckResourceAttr("porkbun_dns_record.test", "content", "192.0.2.1"),
					resource.TestCheckResourceAttr("porkbun_dns_record.test", "ttl", "600"),
					resource.TestCheckResourceAttrSet("porkbun_dns_record.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "porkbun_dns_record.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: importStateIdFunc("porkbun_dns_record.test"),
			},
			// Update testing - change IP
			{
				Config: testAccDNSRecordResourceConfig_A("tftest", "192.0.2.2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("porkbun_dns_record.test", "content", "192.0.2.2"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccDNSRecordResource_AAAA(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccDNSRecordResourceConfig_AAAA("tftest-ipv6", "2001:db8::1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("porkbun_dns_record.test_aaaa", "domain", testDomain),
					resource.TestCheckResourceAttr("porkbun_dns_record.test_aaaa", "name", "tftest-ipv6"),
					resource.TestCheckResourceAttr("porkbun_dns_record.test_aaaa", "type", "AAAA"),
					resource.TestCheckResourceAttr("porkbun_dns_record.test_aaaa", "content", "2001:db8::1"),
					resource.TestCheckResourceAttrSet("porkbun_dns_record.test_aaaa", "id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccDNSRecordResource_CNAME(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccDNSRecordResourceConfig_CNAME("tftest-cname", "example.com"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("porkbun_dns_record.test_cname", "domain", testDomain),
					resource.TestCheckResourceAttr("porkbun_dns_record.test_cname", "name", "tftest-cname"),
					resource.TestCheckResourceAttr("porkbun_dns_record.test_cname", "type", "CNAME"),
					resource.TestCheckResourceAttr("porkbun_dns_record.test_cname", "content", "example.com"),
					resource.TestCheckResourceAttrSet("porkbun_dns_record.test_cname", "id"),
				),
			},
			// Update testing - change target
			{
				Config: testAccDNSRecordResourceConfig_CNAME("tftest-cname", "example.org"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("porkbun_dns_record.test_cname", "content", "example.org"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccDNSRecordResource_TXT(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccDNSRecordResourceConfig_TXT("tftest-txt", "v=test1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("porkbun_dns_record.test_txt", "domain", testDomain),
					resource.TestCheckResourceAttr("porkbun_dns_record.test_txt", "name", "tftest-txt"),
					resource.TestCheckResourceAttr("porkbun_dns_record.test_txt", "type", "TXT"),
					resource.TestCheckResourceAttr("porkbun_dns_record.test_txt", "content", "v=test1"),
					resource.TestCheckResourceAttrSet("porkbun_dns_record.test_txt", "id"),
				),
			},
			// Update testing - change value
			{
				Config: testAccDNSRecordResourceConfig_TXT("tftest-txt", "v=test2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("porkbun_dns_record.test_txt", "content", "v=test2"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccDNSRecordResource_MX(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccDNSRecordResourceConfig_MX("tftest-mx", "mail.example.com", "10"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("porkbun_dns_record.test_mx", "domain", testDomain),
					resource.TestCheckResourceAttr("porkbun_dns_record.test_mx", "name", "tftest-mx"),
					resource.TestCheckResourceAttr("porkbun_dns_record.test_mx", "type", "MX"),
					resource.TestCheckResourceAttr("porkbun_dns_record.test_mx", "content", "mail.example.com"),
					resource.TestCheckResourceAttr("porkbun_dns_record.test_mx", "prio", "10"),
					resource.TestCheckResourceAttrSet("porkbun_dns_record.test_mx", "id"),
				),
			},
			// Update testing - change priority
			{
				Config: testAccDNSRecordResourceConfig_MX("tftest-mx", "mail.example.com", "20"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("porkbun_dns_record.test_mx", "prio", "20"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccDNSRecordResource_RootDomain(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create TXT record on root domain
			{
				Config: testAccDNSRecordResourceConfig_RootTXT("tf-acc-test-root"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("porkbun_dns_record.test_root", "domain", testDomain),
					resource.TestCheckResourceAttr("porkbun_dns_record.test_root", "name", ""),
					resource.TestCheckResourceAttr("porkbun_dns_record.test_root", "type", "TXT"),
					resource.TestCheckResourceAttr("porkbun_dns_record.test_root", "content", "tf-acc-test-root"),
					resource.TestCheckResourceAttrSet("porkbun_dns_record.test_root", "id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccDNSRecordResource_WithNotes(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with notes
			{
				Config: testAccDNSRecordResourceConfig_WithNotes("tftest-notes", "192.0.2.100", "Test record created by Terraform"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("porkbun_dns_record.test_notes", "domain", testDomain),
					resource.TestCheckResourceAttr("porkbun_dns_record.test_notes", "name", "tftest-notes"),
					resource.TestCheckResourceAttr("porkbun_dns_record.test_notes", "type", "A"),
					resource.TestCheckResourceAttr("porkbun_dns_record.test_notes", "content", "192.0.2.100"),
					resource.TestCheckResourceAttr("porkbun_dns_record.test_notes", "notes", "Test record created by Terraform"),
					resource.TestCheckResourceAttrSet("porkbun_dns_record.test_notes", "id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccDNSRecordResource_CustomTTL(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with custom TTL
			{
				Config: testAccDNSRecordResourceConfig_CustomTTL("tftest-ttl", "192.0.2.200", "3600"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("porkbun_dns_record.test_ttl", "domain", testDomain),
					resource.TestCheckResourceAttr("porkbun_dns_record.test_ttl", "name", "tftest-ttl"),
					resource.TestCheckResourceAttr("porkbun_dns_record.test_ttl", "type", "A"),
					resource.TestCheckResourceAttr("porkbun_dns_record.test_ttl", "content", "192.0.2.200"),
					resource.TestCheckResourceAttr("porkbun_dns_record.test_ttl", "ttl", "3600"),
					resource.TestCheckResourceAttrSet("porkbun_dns_record.test_ttl", "id"),
				),
			},
			// Update TTL
			{
				Config: testAccDNSRecordResourceConfig_CustomTTL("tftest-ttl", "192.0.2.200", "7200"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("porkbun_dns_record.test_ttl", "ttl", "7200"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

// Config helper functions

func testAccDNSRecordResourceConfig_A(name, ip string) string {
	return fmt.Sprintf(`
resource "porkbun_dns_record" "test" {
  domain  = %[1]q
  name    = %[2]q
  type    = "A"
  content = %[3]q
  ttl     = "600"
}
`, testDomain, name, ip)
}

func testAccDNSRecordResourceConfig_AAAA(name, ip string) string {
	return fmt.Sprintf(`
resource "porkbun_dns_record" "test_aaaa" {
  domain  = %[1]q
  name    = %[2]q
  type    = "AAAA"
  content = %[3]q
  ttl     = "600"
}
`, testDomain, name, ip)
}

func testAccDNSRecordResourceConfig_CNAME(name, target string) string {
	return fmt.Sprintf(`
resource "porkbun_dns_record" "test_cname" {
  domain  = %[1]q
  name    = %[2]q
  type    = "CNAME"
  content = %[3]q
  ttl     = "600"
}
`, testDomain, name, target)
}

func testAccDNSRecordResourceConfig_TXT(name, value string) string {
	return fmt.Sprintf(`
resource "porkbun_dns_record" "test_txt" {
  domain  = %[1]q
  name    = %[2]q
  type    = "TXT"
  content = %[3]q
  ttl     = "600"
}
`, testDomain, name, value)
}

func testAccDNSRecordResourceConfig_MX(name, target, priority string) string {
	return fmt.Sprintf(`
resource "porkbun_dns_record" "test_mx" {
  domain  = %[1]q
  name    = %[2]q
  type    = "MX"
  content = %[3]q
  prio    = %[4]q
  ttl     = "600"
}
`, testDomain, name, target, priority)
}

func testAccDNSRecordResourceConfig_RootTXT(value string) string {
	return fmt.Sprintf(`
resource "porkbun_dns_record" "test_root" {
  domain  = %[1]q
  name    = ""
  type    = "TXT"
  content = %[2]q
  ttl     = "600"
}
`, testDomain, value)
}

func testAccDNSRecordResourceConfig_WithNotes(name, ip, notes string) string {
	return fmt.Sprintf(`
resource "porkbun_dns_record" "test_notes" {
  domain  = %[1]q
  name    = %[2]q
  type    = "A"
  content = %[3]q
  ttl     = "600"
  notes   = %[4]q
}
`, testDomain, name, ip, notes)
}

func testAccDNSRecordResourceConfig_CustomTTL(name, ip, ttl string) string {
	return fmt.Sprintf(`
resource "porkbun_dns_record" "test_ttl" {
  domain  = %[1]q
  name    = %[2]q
  type    = "A"
  content = %[3]q
  ttl     = %[4]q
}
`, testDomain, name, ip, ttl)
}
