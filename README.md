# Terraform LDAP [![CircleCI](https://circleci.com/gh/Pryz/terraform-provider-ldap.svg?style=svg)](https://circleci.com/gh/Pryz/terraform-provider-ldap)

## Installation

You can easily install the latest version with the following :

```
go get -u github.com/Pryz/terraform-provider-ldap
```

Then add the plugin to your local `.terraformrc` :

```
cat >> ~/.terraformrc <<EOF
providers {
  ldap = "${GOPATH}/bin/terraform-provider-ldap"
}
EOF
```

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

## Resource LDAP Object example

```
resource "ldap_object" "foo" {
  dn = "uid=foo"
  base_dn = "dc=example,dc=com"

  object_classes = [
    "inetOrgPerson",
    "posixAccount",
  ]

  attribute {
    name = "sn"
    value = "10"
  }
  attribute {
    name = "cn"
    value = "bar"
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
    value = "/home/billy"
  }
  attribute {
    name = "loginShell"
    value = "/bin/bash"
  }

}
```

Of course the Bind User will need write access.

## Limitations

Currently this provider doesn't handle updates. To change records it will delete and create a new one.
I don't see any problem by doing that for the moment. If you do feel free to create me an Issue.
