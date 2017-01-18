provider "ldap" {
  ldap_host = "localhost"
  ldap_port = 389
  use_tls = false
  bind_user = "cn=admin,dc=example,dc=com"
  bind_password = "admin"
}

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
