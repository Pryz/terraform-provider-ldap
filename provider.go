package ldap

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"ldap_host": &schema.Schema{
				Type: schema.TypeString,
				Required: true,
				DefaultFunc: schema.EnvDefaultFunc("LDAP_HOST", nil),
				Description: descriptions["ldap_host"],
			},
			"ldap_port": &schema.Schema{
				Type: schema.TypeInt,
				Required: false,
				DefaultFunc: schema.EnvDefaultFunc("LDAP_PORT", 389),
				Description: descriptions["ldap_port"],
			},
			"use_tls": &schema.Schema{
				Type: schema.TypeBool,
				Required: false,
				DefaultFunc: schema.EnvDefaultFunc("LDAP_USE_TLS", true),
				Description: descriptions["ldap_port"],
			},
			"bind_user": &schema.Schema{
				Type: schema.TypeString,
				Required: true,
				DefaultFunc: schema.EnvDefaultFunc("LDAP_BIND_USER", nil),
				Description: descriptions["bind_user"],
			},
			"bind_password": &schema.Schema{
				Type: schema.TypeString,
				Required: true,
				DefaultFunc: schema.EnvDefaultFunc("LDAP_BIND_PASSWORD", nil),
				Description: descriptions["bind_password"],
			},
		},

		ResourceMap: map[string]*schema.Resource{

		},
		ConfigureFunc: configureProvider,
	}
}

func init() {
	descriptions = map[string]string{
		"ldap_host": "The LDAP host to initiate the conneciton.",

		"ldap_port": "The LDAP port to initiate the conneciton. Default : 389.",

		"use_tls": "Use TLS to secure the connection. Default: true.",

		"bind_user": "Bind user to be used for the LDAP request.",

		"bind_password": "Password to authenticate the Bind user.",
	}
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		LdapHost: d.Get("ldap_host").(string),
		LdapPort: d.Get("ldap_port").(int),
		UseTLS: d.Get("use_tls").(bool),
		BindUser: d.Get("bind_user").(string),
		BindPassword: d.Get("bind_password").(string),
	}

	if conn, err := config.initiateAndBind(); err != nil {
		return nil, err
	}

	return conn, nil
}
