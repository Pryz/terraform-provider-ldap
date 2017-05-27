# Terraform LDAP 

[![CircleCI](https://circleci.com/gh/Pryz/terraform-provider-ldap.svg?style=svg)](https://circleci.com/gh/Pryz/terraform-provider-ldap)

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
    ldap_host = "ldap.example.org"
    ldap_port = 389
    use_tls = true
    bind_user = "cn=admin,dc=example,dc=com"
    bind_password = "admin"
}
```

## Resource LDAP Object example

```
resource "ldap_object" "foo" {
    # DN must be complete (no RDN!)
    dn = "uid=foo,dc=example,dc=com"

    # classes are specified as an array
    object_classes = [
        "inetOrgPerson",
        "posixAccount",
    ]

    # attributes are sepcified as a set of 1-element maps
    attributes = [
        { sn              = "10" },
        { cn              = "bar" },
        { uidNumber       = "1234" },
        { gidNumber       = "1234" },
        { homeDirectory   = "/home/billy" },
        { loginShell      = "/bin/bash" },
        # when an attribute has multiple values, it must be specified multiple times
        { mail            = "billy@example.com" },
        { mail            = "admin@example.com" },
    ]
}
```

The Bind User must have write access for resource creation to succeed.

## Features

This provider is feature complete; it supports resource creation, reading, update 
and deletion; it can be used to create nested resources at all levels of the
hierarchy, provided the proper (implicit or explicit) dependencies are declared.
When it comes to updating an object, the plugin will calculate the set of 
attributes that need to be added, modified and removed and will surgically 
operate on the remote object.

## Limitations

This provider supports TLS, but certificate verification is not enabled yet; all
connections are through TCP, no UDP support yet.