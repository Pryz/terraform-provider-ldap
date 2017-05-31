package main

import (
	"errors"
	"fmt"
	"testing"

	ldap "gopkg.in/ldap.v2"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccLDAPObject_Basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLDAPObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLDAPObjectConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLDAPObjectExists("ldap_object.jdoe"),
					resource.TestCheckResourceAttr("ldap_object.jdoe", "dn", "uid=jdoe,dv=example,dc=com"),
					//resource.TestCheckResourceAttr("ldap_object.jdoe", "base_dn", "dc=example,dc=com"),
					resource.TestCheckResourceAttr("ldap_object.jdoe", "object_classes.0", "inetOrgPerson"),
					resource.TestCheckResourceAttr("ldap_object.jdoe", "object_classes.1", "posixAccount"),
					testAccCheckLDAPObjectAttributes("ldap_object.jdoe"),
				),
			},
		},
	})
}

func testAccCheckLDAPObjectDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ldap.Conn)
	for _, r := range s.RootModule().Resources {
		dn := r.Primary.Attributes["dn"]
		sr, err := helperSearchRequest(dn, conn)
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

func testAccCheckLDAPObjectExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*ldap.Conn)
		for _, r := range s.RootModule().Resources {
			dn := r.Primary.Attributes["dn"]
			sr, err := helperSearchRequest(dn, conn)
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

func testAccCheckLDAPObjectAttributes(n string) resource.TestCheckFunc {
	//TODO: check sanity of the ldap object attributes
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("LDAP Object not found: %s", n)
		}
		return nil
	}
}

func helperSearchRequest(dn string, conn *ldap.Conn) (*ldap.SearchResult, error) {

	// search by primary key (that is, set the DN as base DN and use a "base
	// object" scope); no attributes are retrieved since we are onòy checking
	// for existence; all objects have an "objectClass" attribute, so the filter
	// is a "match all"
	request := ldap.NewSearchRequest(
		dn,
		ldap.ScopeBaseObject,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		"(objectClass=*)",
		[]string{"*"},
		nil,
	)

	return conn.Search(request)
}

const testAccCheckLDAPObjectConfig = `
resource "ldap_object" "jdoe" {
  dn = "uid=jdoe,dc=example,dc=com"
	
  object_classes    = [
    "inetOrgPerson",
    "posixAccount",
  ]

  attributes        = [
		{ sn            = "Doe"}, 		
		{ givenName		  = "John"},
		{ cn			      = "John Doe"},
		{ displayName	  = "Mr. John K. Doe, esq."},
		{ mail 			    = "john.doe@example.com" },
		{ mail			    = "a123456@internal.example.com" },
		{ mail			    = "jdoe@example.com" },
		{ userPassword  = "password" },
		{ description	  = "The best programmer in the world." },
  	{ uidNumber     = "1234" },
   	{ gidNumber     = "1234" },
   	{ homeDirectory = "/home/jdoe"},
   	{ loginShell    = "/bin/bash" }
	]
}
`
