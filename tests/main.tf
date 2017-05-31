/*
 * LDAP is the OpenLDAP server.
 */
provider "ldap" {
	ldap_host           = "localhost"
	ldap_port 		   	= 389
	use_tls 		    = false
	bind_user 		   	= "cn=admin,dc=example,dc=com"
	bind_password       = "admin"
}

/*
 * Set up an OU to contain users first.
 *
 * NOTE: objects are identified by DN; BaseDN is not necessary since the
 * Terraform LDAP Provider looks up objects by primary key (DN); attributes are
 * optional, if the object's class does not require them. 
 */ 
resource "ldap_object" "users_example_com" {
	dn 					= "ou=users,dc=example,dc=com"
	object_classes 	    = [ "top", "organizationalUnit" ]
# 	object_classes      = [ "organizationalUnit" ]
}

/*
 * Set up one or more users inside the users' OU.
 *
 * NOTE: objects can be nested; to nest an object and to establish implicit 
 * dependency between containee and container, use the parent's DN and the
 * new object's RDN to provide the new object's DN (DN=RDN,ParentDN).
 */
resource "ldap_object" "a123456" {
	dn                  = "uid=a123456,${ldap_object.users_example_com.dn}"		
	object_classes      = ["inetOrgPerson", "posixAccount"]
	attributes          = [		
		{ sn            = "Doe"}, 		
		{ givenName		= "John"},
		{ cn			= "John Doe"},
		{ displayName	= "Mr. John K. Doe, esq."},
		{ mail 			= "john.doe@example.com" },
#		{ mail			= "a123456@internal.example.com" },
		{ mail			= "jdoe@example.com" },
		{ userPassword  = "password" },
#		{ description	= "The best programmer in the world." },
    	{ uidNumber     = "1234" },
    	{ gidNumber     = "1234" },
    	{ homeDirectory = "/home/jdoe"},
    	{ loginShell    = "/bin/bash" }
	]
}
