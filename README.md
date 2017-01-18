# [WIP] Terraform LDAP

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

Of course the Bind User will need write access.

## Limitations

Currently this provider doesn't handle updates. To change records it will delete and create a new one.
I don't see any problem by doing that for the moment. If you do feel free to create me an Issue.
