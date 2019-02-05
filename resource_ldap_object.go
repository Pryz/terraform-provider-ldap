package main

import (
	"bytes"
	"log"

	"fmt"
	"strings"

	"os"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"gopkg.in/ldap.v2"
)

func resourceLDAPObject() *schema.Resource {
	return &schema.Resource{
		Create: resourceLDAPObjectCreate,
		Read:   resourceLDAPObjectRead,
		Update: resourceLDAPObjectUpdate,
		Delete: resourceLDAPObjectDelete,
		Exists: resourceLDAPObjectExists,

		Importer: &schema.ResourceImporter{
			State: resourceLDAPObjectImport,
		},

		Schema: map[string]*schema.Schema{
			"dn": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The Distinguished Name (DN) of the object, as the concatenation of its RDN (unique among siblings) and its parent's DN.",
				Required:    true,
				ForceNew:    true,
			},
			"object_classes": &schema.Schema{
				Type:        schema.TypeSet,
				Description: "The set of classes this object conforms to (e.g. organizationalUnit, inetOrgPerson).",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Required:    true,
			},
			"attributes": &schema.Schema{
				Type:        schema.TypeSet,
				Description: "The map of attributes of this object; each attribute can be multi-valued.",
				Set:         attributeHash,
				MinItems:    0,

				Elem: &schema.Schema{
					Type:        schema.TypeMap,
					Description: "The list of values for a given attribute.",
					MinItems:    1,
					MaxItems:    1,
					Elem: &schema.Schema{
						Type:        schema.TypeString,
						Description: "The individual value for the given attribute.",
					},
				},
				Optional: true,
			},
		},
	}
}

func resourceLDAPObjectImport(d *schema.ResourceData, meta interface{}) (imported []*schema.ResourceData, err error) {
	d.Set("dn", d.Id())
	err = readLDAPObjectImpl(d, meta, false)
	if path := os.Getenv("TF_LDAP_IMPORTER_PATH"); path != "" {
		log.Printf("[DEBUG] ldap_object::import - dumping imported object to %q", path)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			// the export file does not exist
			if file, err := os.Create(path); err == nil {
				defer file.Close()
				id := d.Id()
				tokens := strings.Split(id, ",")
				if len(tokens) > 0 {
					tokens = strings.Split(tokens[0], "=")
					if len(tokens) >= 1 {
						id = tokens[1]
						//resource "ldap_object" "a123456" {
						file.WriteString(fmt.Sprintf("resource \"ldap_object\" %q {\n", id))
						//	dn = "uid=a123456,dc=example,dc=com"
						file.WriteString(fmt.Sprintf("  dn = %q\n", d.Id()))
						//  object_classes = ["inetOrgPerson", "posixAccount"]
						classes := []string{}
						for _, class := range d.Get("object_classes").(*schema.Set).List() {
							//classes[i] = fmt.Sprintf("\"%s\"", class)
							classes = append(classes, fmt.Sprintf("%q", class))
						}
						file.WriteString(fmt.Sprintf("  object_classes = [ %s ]\n", strings.Join(classes, ", ")))
						if attributes, ok := d.GetOk("attributes"); ok {
							attributes := attributes.(*schema.Set).List()
							if len(attributes) > 0 {
								//  attributes = [
								file.WriteString("  attributes = [\n")
								for _, attribute := range attributes {
									for name, value := range attribute.(map[string]interface{}) {
										//    { sn = "Doe" },
										file.WriteString(fmt.Sprintf("    { %s = %q },\n", name, value.(string)))
									}
								}
								// ]
								file.WriteString("  ]\n")
							}
						}
						file.WriteString("}\n")
					}
				}
			}
		}
	}
	imported = append(imported, d)
	return
}

func resourceLDAPObjectExists(d *schema.ResourceData, meta interface{}) (b bool, e error) {
	config := meta.(*Config)
	conn, err := config.initiateAndBind()
	if err != nil {
		return false, err
	}
	dn := d.Get("dn").(string)

	log.Printf("[DEBUG] ldap_object::exists - checking if %q exists", dn)

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
		nil,
		nil,
	)

	var _ *ldap.SearchResult
	_, err = conn.Search(request)
	if err != nil {
		if err, ok := err.(*ldap.Error); ok {
			if err.ResultCode == 32 { // no such object
				log.Printf("[WARN] ldap_object::exists - lookup for %q returned no value: deleted on server?", dn)
				return false, nil
			}
		}
		log.Printf("[DEBUG] ldap_object::exists - lookup for %q returned an error %v", dn, err)
		return false, err
	}

	log.Printf("[DEBUG] ldap_object::exists - object %q exists", dn)
	return true, nil
}

func resourceLDAPObjectCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.initiateAndBind()
	if err != nil {
		return err
	}
	dn := d.Get("dn").(string)

	log.Printf("[DEBUG] ldap_object::create - creating a new object under %q", dn)

	request := ldap.NewAddRequest(dn)

	// retrieve classe from HCL
	objectClasses := []string{}
	for _, oc := range (d.Get("object_classes").(*schema.Set)).List() {
		log.Printf("[DEBUG] ldap_object::create - object %q has class: %q", dn, oc.(string))
		objectClasses = append(objectClasses, oc.(string))
	}
	request.Attribute("objectClass", objectClasses)

	// if there is a non empty list of attributes, loop though it and
	// create a new map collecting attribute names and its value(s); we need to
	// do this because we could not model the attributes as a map[string][]string
	// due to an appareent limitation in HCL; we have a []map[string]string, so
	// we loop through the list and accumulate values when they share the same
	// key, then we use these as attributes in the LDAP client.
	if v, ok := d.GetOk("attributes"); ok {
		attributes := v.(*schema.Set).List()
		if len(attributes) > 0 {
			log.Printf("[DEBUG] ldap_object::create - object %q has %d attributes", dn, len(attributes))
			m := make(map[string][]string)
			for _, attribute := range attributes {
				log.Printf("[DEBUG] ldap_object::create - %q has attribute of type %T", dn, attribute)
				// each map should only have one entry (see resource declaration)
				for name, value := range attribute.(map[string]interface{}) {
					log.Printf("[DEBUG] ldap_object::create - %q has attribute[%v] => %v (%T)", dn, name, value, value)
					m[name] = append(m[name], value.(string))
				}
			}
			// now loop through the map and add attributes with theys value(s)
			for name, values := range m {
				request.Attribute(name, values)
			}
		}
	}

	err = client.Add(request)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] ldap_object::create - object %q added to LDAP server", dn)

	d.SetId(dn)
	return resourceLDAPObjectRead(d, meta)
}

func resourceLDAPObjectRead(d *schema.ResourceData, meta interface{}) error {
	return readLDAPObjectImpl(d, meta, true)
}

func resourceLDAPObjectUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.initiateAndBind()
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] ldap_object::update - performing update on %q", d.Id())

	request := ldap.NewModifyRequest(d.Id())

	// handle objectClasses
	if d.HasChange("object_classes") {
		classes := []string{}
		for _, oc := range (d.Get("object_classes").(*schema.Set)).List() {
			classes = append(classes, oc.(string))
		}
		log.Printf("[DEBUG] ldap_object::update - updating classes of %q, new value: %v", d.Id(), classes)
		request.ReplaceAttributes = []ldap.PartialAttribute{
			ldap.PartialAttribute{
				Type: "objectClass",
				Vals: classes,
			},
		}
	}

	if d.HasChange("attributes") {

		o, n := d.GetChange("attributes")
		log.Printf("[DEBUG] ldap_object::update - \n%s", printAttributes("old attributes map", o))
		log.Printf("[DEBUG] ldap_object::update - \n%s", printAttributes("new attributes map", n))

		added, changed, removed := computeDeltas(o.(*schema.Set), n.(*schema.Set))
		if len(added) > 0 {
			log.Printf("[DEBUG] ldap_object::update - %d attributes added", len(added))
			request.AddAttributes = added
		}
		if len(changed) > 0 {
			log.Printf("[DEBUG] ldap_object::update - %d attributes changed", len(changed))
			if request.ReplaceAttributes == nil {
				request.ReplaceAttributes = changed
			} else {
				request.ReplaceAttributes = append(request.ReplaceAttributes, changed...)
			}
		}
		if len(removed) > 0 {
			log.Printf("[DEBUG] ldap_object::update - %d attributes removed", len(removed))
			request.DeleteAttributes = removed
		}
	}

	err = client.Modify(request)
	if err != nil {
		log.Printf("[ERROR] ldap_object::update - error modifying LDAP object %q with values %v", d.Id(), err)
		return err
	}
	return resourceLDAPObjectRead(d, meta)
}

func resourceLDAPObjectDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.initiateAndBind()
	if err != nil {
		return err
	}
	dn := d.Get("dn").(string)

	log.Printf("[DEBUG] ldap_object::delete - removing %q", dn)

	request := ldap.NewDelRequest(dn, nil)

	err = client.Del(request)
	if err != nil {
		log.Printf("[ERROR] ldap_object::delete - error removing %q: %v", dn, err)
		return err
	}
	log.Printf("[DEBUG] ldap_object::delete - %q removed", dn)
	return nil
}

func readLDAPObjectImpl(d *schema.ResourceData, meta interface{}, updateState bool) error {
	config := meta.(*Config)
	client, err := config.initiateAndBind()
	if err != nil {
		return err
	}
	dn := d.Get("dn").(string)

	log.Printf("[DEBUG] ldap_object::read - looking for object %q", dn)

	// when searching by DN, you don't need t specify the base DN a search
	// filter a "subtree" scope: just put the DN (i.e. the primary key) as the
	// base DN with a "base object" scope, and the returned object will be the
	// entry, if it exists
	request := ldap.NewSearchRequest(
		dn,
		ldap.ScopeBaseObject,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		"(objectclass=*)",
		[]string{"*"},
		nil,
	)

	sr, err := client.Search(request)
	if err != nil {
		if err, ok := err.(*ldap.Error); ok {
			if err.ResultCode == 32 && updateState { // no such object
				log.Printf("[WARN] ldap_object::read - object not found, removing %q from state because it no longer exists in LDAP", dn)
				d.SetId("")
				return nil
			}
		}
		log.Printf("[DEBUG] ldap_object::read - lookup for %q returned an error %v", dn, err)
		return err
	}

	log.Printf("[DEBUG] ldap_object::read - query for %q returned %v", dn, sr)

	d.SetId(dn)
	d.Set("object_classes", sr.Entries[0].GetAttributeValues("objectClass"))

	// now deal with attributes
	set := &schema.Set{
		F: attributeHash,
	}

	for _, attribute := range sr.Entries[0].Attributes {
		log.Printf("[DEBUG] ldap_object::read - treating attribute %q of %q (%d values: %v)", attribute.Name, dn, len(attribute.Values), attribute.Values)
		if attribute.Name == "objectClass" {
			// skip: we don't treat object classes as ordinary attributes
			log.Printf("[DEBUG] ldap_object::read - skipping attribute %q of %q", attribute.Name, dn)
			continue
		}
		log.Printf("[DEBUG] ldap_object::read - adding attribute %q to %q (%d values)", attribute.Name, dn, len(attribute.Values))
		// now add each value as an individual entry into the object, because
		// we do not handle name => []values, and we have a set of maps each
		// holding a single entry name => value; multiple maps may share the
		// same key.
		for _, value := range attribute.Values {
			log.Printf("[DEBUG] ldap_object::read - for %q, setting %q => %q", dn, attribute.Name, value)
			set.Add(map[string]interface{}{
				attribute.Name: value,
			})
		}
	}

	if err := d.Set("attributes", set); err != nil {
		log.Printf("[WARN] ldap_object::read - error setting LDAP attributes for %q : %v", dn, err)
		return err
	}
	return nil
}

// computes the hash of the map representing an attribute in the attributes set
func attributeHash(v interface{}) int {
	m := v.(map[string]interface{})
	var buffer bytes.Buffer
	buffer.WriteString("map {")
	for k, v := range m {
		buffer.WriteString(fmt.Sprintf("%q := %q;", k, v.(string)))
	}
	buffer.WriteRune('}')
	text := buffer.String()
	hash := hashcode.String(text)
	return hash
}

func printAttributes(prefix string, attributes interface{}) string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("%s: {\n", prefix))
	if attributes, ok := attributes.(*schema.Set); ok {
		for _, attribute := range attributes.List() {
			for k, v := range attribute.(map[string]interface{}) {
				buffer.WriteString(fmt.Sprintf("    %q: %q\n", k, v.(string)))
			}
		}
		buffer.WriteRune('}')
	}
	return buffer.String()
}

func computeDeltas(os, ns *schema.Set) (added, changed, removed []ldap.PartialAttribute) {

	rk := NewSet() // names of removed attributes
	for _, v := range os.Difference(ns).List() {
		for k := range v.(map[string]interface{}) {
			rk.Add(k)
		}
	}

	ak := NewSet() // names of added attributes
	for _, v := range ns.Difference(os).List() {
		for k := range v.(map[string]interface{}) {
			ak.Add(k)
		}
	}

	kk := NewSet() // names of kept attributes
	for _, v := range ns.Intersection(os).List() {
		for k := range v.(map[string]interface{}) {
			kk.Add(k)
		}
	}

	ck := NewSet() // names of changed attributes

	// loop over remove attributes' names
	for _, k := range rk.List() {
		if !ak.Contains(k) && !kk.Contains(k) {
			// one value under this name has been removed, no other value has
			// been added back, and there is no further value under the same
			// name among those that were untouched; this means that it has
			// been dropped and must go among the RemovedAttributes
			log.Printf("[DEBUG} ldap_object::deltas - dropping attribute %q", k)
			removed = append(removed, ldap.PartialAttribute{
				Type: k,
				Vals: []string{},
			})
		} else {
			ck.Add(k)
		}
	}

	for _, k := range ak.List() {
		if !rk.Contains(k) && !kk.Contains(k) {
			// this is the first value under this name: no value is being
			// removed and no value is being kept; so we're adding this new
			// attribute to the LDAP object (AddedAttributes), getting all
			// the values under this name from the new set
			values := []string{}
			for _, m := range ns.List() {
				for mk, mv := range m.(map[string]interface{}) {
					if k == mk {
						values = append(values, mv.(string))
					}
				}
			}
			added = append(added, ldap.PartialAttribute{
				Type: k,
				Vals: values,
			})
			log.Printf("[DEBUG} ldap_object::deltas - adding new attribute %q with values %v", k, values)
		} else {
			ck.Add(k)
		}
	}

	// now loop over changed attributes and
	for _, k := range ck.List() {
		// the attributes in this set have been changed, in that a new value has
		// been added or removed and it was not the last/first one; so we're
		// adding this new attribute to the LDAP object (ModifiedAttributes),
		// getting all the values under this name from the new set
		values := []string{}
		for _, m := range ns.List() {
			for mk, mv := range m.(map[string]interface{}) {
				if k == mk {
					values = append(values, mv.(string))
				}
			}
		}
		changed = append(added, ldap.PartialAttribute{
			Type: k,
			Vals: values,
		})
		log.Printf("[DEBUG} ldap_object::deltas - changing attribute %q with values %v", k, values)
	}
	return
}
