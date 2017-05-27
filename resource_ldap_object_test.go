package main

import (
	"errors"
	"fmt"
	"testing"

	ldap "gopkg.in/ldap.v2"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccLdapObject_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLdapObjectDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckLdapObjectConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLdapObjectExists("ldap_object.foo"),
					resource.TestCheckResourceAttr(
						"ldap_object.foo", "dn", "uid=foo"),
					resource.TestCheckResourceAttr(
						"ldap_object.foo", "base_dn", "dc=example,dc=com"),
					resource.TestCheckResourceAttr(
						"ldap_object.foo", "object_classes.0", "inetOrgPerson"),
					resource.TestCheckResourceAttr(
						"ldap_object.foo", "object_classes.1", "posixAccount"),
					testAccCheckLdapObjectAttributes("ldap_object.foo"),
				),
			},
		},
	})
}

func testAccCheckLdapObjectDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ldap.Conn)
	for _, r := range s.RootModule().Resources {
		dn := r.Primary.Attributes["dn"]
		baseDN := r.Primary.Attributes["base_dn"]
		sr, err := helperSearchRequest(dn, baseDN, conn)
		if err != nil {
			return err
		}
		if len(sr.Entries) != 0 {
			err = errors.New("Number of records greater than 0 for " + dn)
			return err
		}
	}
	return nil
}

func testAccCheckLdapObjectExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*ldap.Conn)
		for _, r := range s.RootModule().Resources {
			dn := r.Primary.Attributes["dn"]
			baseDn := r.Primary.Attributes["base_dn"]
			sr, err := helperSearchRequest(dn, baseDn, conn)
			if err != nil {
				return err
			}

			if len(sr.Entries) != 1 {
				err = errors.New("Number of record different than 1 for " + dn)
				return err
			}
		}
		return nil
	}
}

func testAccCheckLdapObjectAttributes(n string) resource.TestCheckFunc {
	//TODO: check sanity of the ldap object attributes
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Ldap Object not found: %s", n)
		}
		return nil
	}
}

const testAccCheckLdapObjectConfig = `
resource "ldap_object" "foo" {
  dn = "uid=foo,dc=example,dc=com"
  object_classes = [
    "inetOrgPerson",
    "posixAccount",
  ]

  attributes = [
		{ sn = "10"
  }
  attribute {
    name = "cn"
    value = "foo"
  }
  attribute {
    name = "uidNumber"
    value = "1234"
  }
  attribute {
    name = "gidNumber"
    value = "1234"
  }
  attribute {
    name = "homeDirectory"
    value = "/home/foo"
  }
  attribute {
    name = "loginShell"
    value = "/bin/bash"
  }

}
`

func helperSearchRequest(dn string, baseDn string, conn *ldap.Conn) (*ldap.SearchResult, error) {
	searchRequest := ldap.NewSearchRequest(baseDn, ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false, fmt.Sprintf("(%s)", dn), []string{"*"}, nil)
	return conn.Search(searchRequest)
}
