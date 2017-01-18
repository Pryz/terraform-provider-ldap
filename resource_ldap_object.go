package main

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/hashicorp/terraform/helper/schema"
	"gopkg.in/ldap.v2"
)

func resourceLdapObject() *schema.Resource {
	return &schema.Resource{
		Create: resourceLdapObjectCreate,
		Read:   resourceLdapObjectRead,
		Update: resourceLdapObjectUpdate,
		Delete: resourceLdapObjectDelete,
		Exists: resourceLdapObjectExists,

		Schema: map[string]*schema.Schema{
			"dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"base_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"object_classes": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"attribute": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"value": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceLdapObjectExists(d *schema.ResourceData, meta interface{}) (b bool, e error) {
	conn := meta.(*ldap.Conn)
	dn := d.Get("dn").(string)
	base_dn := d.Get("base_dn").(string)

	searchRequest := ldap.NewSearchRequest(
		base_dn,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(%s)", dn),
		nil,
		nil,
	)

	sr, err := conn.Search(searchRequest)
	if err != nil {
		return false, err
	}

	if len(sr.Entries) == 0 {
		return false, nil
	} else if len(sr.Entries) > 1 {
		err = errors.New("More than one record found for " + dn)
		return false, err
	}

	return true, nil
}

func resourceLdapObjectCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*ldap.Conn)
	dn := d.Get("dn").(string)
	base_dn := d.Get("base_dn").(string)

	addRequest := ldap.NewAddRequest(fmt.Sprintf("%s,%s", dn, base_dn))

	objectClasses := []string{}
	for _, oc := range d.Get("object_classes").([]interface{}) {
		objectClasses = append(objectClasses, oc.(string))
	}
	addRequest.Attribute("objectClass", objectClasses)

	if attributes := d.Get("attribute").(*schema.Set); attributes.Len() > 0 {
		for _, attr := range attributes.List() {
			m := attr.(map[string]interface{})
			addRequest.Attribute(m["name"].(string), []string{m["value"].(string)})
		}
	}

	err := conn.Add(addRequest)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%d", rand.Int()))
	return nil
}

func resourceLdapObjectRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*ldap.Conn)
	dn := d.Get("dn").(string)
	base_dn := d.Get("base_dn").(string)

	searchRequest := ldap.NewSearchRequest(
		base_dn,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(%s)", dn),
		nil,
		nil,
	)

	sr, err := conn.Search(searchRequest)
	if err != nil {
		return err
	}

	if len(sr.Entries) == 0 {
		err = errors.New("No record found for " + dn)
		return err
	} else if len(sr.Entries) > 1 {
		err = errors.New("More than one record found for " + dn)
		return err
	}
	return nil
}

func resourceLdapObjectUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceLdapObjectDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
