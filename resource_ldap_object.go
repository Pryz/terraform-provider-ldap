package ldap

import (
	"errors"

	"github.com/hashicorp/terraform/helper/schema"
	"gopkg.in/ldap.v2"
)

func resourceLdapObject() *schema.Resource {
	return &schema.Resource{
		Create: resourceLdapObjectCreate,
		Read:		resourceLdapObjectRead,
		Update: resourceLdapObjectUpdate,
		Delete: resourceLdapObjectDelete,
		Exists: resourceLdapObjectExists,
		Importer: nil,

		Schema: map[string]*schema.Schema{
			"dn": &schema.Schema{
				Type:	schema.TypeString,
				Required: true,
			},
			"base_dn": &schema.Schema{
				Type:	schema.TypeString,
				Required: true,
			},
			"object_class": &schema.Schema{
				Type: schema.TypeMap,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"value": &schema.Schema{
							Type: schema.TypeString,
							Optional: false
						},
					},
				},
			},
			"attribute": &schema.Schema{
				Type: schema.TypeMap,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type: schema.TypeString,
							Optional: false
						},
						"value": &schema.Schema{
							Type: schema.TypeString,
							Optional: false
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
		fmt.Sprintf("(&(objectClass=*)&(%s))", dn),
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
		err = errors.New("More than one record found for %s", dn)
		return false, err
	}

	return true, nil
}

func resourceLdapObjectCreate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceLdapObjectRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceLdapObjectUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceLdapObjectDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
