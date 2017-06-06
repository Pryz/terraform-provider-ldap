# Terraform LDAP 

[![CircleCI](https://circleci.com/gh/dihedron/terraform-provider-ldap.svg?style=svg)](https://circleci.com/gh/dihedron/terraform-provider-ldap)

## Note

This Terraform provider is a fork of 
[a previous implementation by Pryz](https://github.com/Pryz/terraform-provider-ldap), which is still available.
The necessity of forking and continuing development on an independent repository
rises from the need of implementing things at a much faster pace and being able 
to commit and have changes available as soon as possible. The initial set of
changes - which included an almost complete rewrite of the provider - were 
contributed back to the upstream repository, but most of new developments will
from now on happen here.

## Installation

You can easily install the latest version with the following :

```
go get -u github.com/dihedron/terraform-provider-ldap
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
    use_tls = false
    bind_user = "cn=admin,dc=example,dc=com"
    bind_password = "admin"
}
```
Note: if you want to use TLS, the LDAP port must be changed accordingly 
(typically, port 636 is used for secure connections).

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

    # attributes are specified as a set of 1-element maps
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

This provider is feature complete.
As of the latest release, it supports resource creation, reading, update, deletion
and importing.
It can be used to create nested resources at all levels of the hierarchy, 
provided the proper (implicit or explicit) dependencies are declared.
When updating an object, the plugin computes the minimum set of attributes that 
need to be added, modified and removed and surgically operates on the remote 
object to bring it up to date.
When importing existing LDAP objects into the Terraform state, the plugin can
automatically generate a .tf file with the relevant information, so that the 
following ```terraform apply``` does not drop the imported resource out of the
remote LDAP server due to it missing in the local ```.tf``` files.
In order to have the plugin generate this file, put the name of the output file
(which must *not* exist on disk) in the ```TF_LDAP_IMPORTER_PATH``` environment 
variable, like this:
```
$> export TF_LDAP_IMPORTER_PATH=a123456.tf 
$> terraform import ldap_object.a123456 uid=a123456,ou=users,dc=example,dc=com
```
and the plugin will create the ```a123456.tf``` file with the proper information.
Then merge this file into your existing ```.tf``` file(s).

## Limitations

This provider supports TLS, but certificate verification is not enabled yet; all
connections are through TCP, no UDP support yet.