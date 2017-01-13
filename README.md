# [WIP] Terraform LDAP

NOT WORKING YET. Come back in few days :)

## Provider example

```
provider "ldap" {
  ldap_host = "ldap.mydomain.com"
  ldap_port = 689
  use_tls = true
  bind_user = "foo"
  bind_password = "bar"
}
```

## DNS example with LDAP object

```
resource "ldap_object" "dns_test_public" {
  dn = "relativeDomainName=hostname42"
  base_dn = "zoneName=${var.public_zone},ou=DNS,dc=mydomain,dc=com"

  object_class {
    value = "dNSZone"
  }

  attribute {
    name = "relativeDomainName"
    value = "hostname42"
  }
  
  attribute {
    name = "zoneName"
    value = "${var.public_zone}"
  }

  attribute {
    name = "aRecord"
    value = "198.98.42.42"
  }

  attribute {
    name = "dNSClass"
    value = "IN"
  }

  attribute {
    name = "dNSTTL"
    value = "3600"
  }
}
```
