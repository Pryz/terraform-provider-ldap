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

## Docker build image

You can build the final binary by calling the ```./build.sh``` which uses
Docker, and copies the resulting binary ```terraform-provider-ldap``` to the
```bin``` directory:

```
$ ./build.sh
================================================================
Building docker image...
================================================================
Sending build context to Docker daemon  172.5kB
Step 1/6 : FROM golang:1.11.3-stretch
 ---> bbf428bade77
Step 2/6 : COPY . /go/src/terraform-provider-ldap
 ---> 97272b52a9b6
Step 3/6 : WORKDIR /go/src/terraform-provider-ldap
 ---> Running in 3935be9d7ea3
Removing intermediate container 3935be9d7ea3
 ---> 72cc7041aec9
Step 4/6 : RUN go get .
 ---> Running in 5c00c77b9f91
Removing intermediate container 5c00c77b9f91
 ---> aeb8bb9c41b9
Step 5/6 : RUN CGO_ENABLED=0 GOOS=linux go install -a -ldflags '-extldflags "-static"' .
 ---> Running in df11c6d06069
Removing intermediate container df11c6d06069
 ---> 61bc6ebab4dd
Step 6/6 : ENTRYPOINT ["/bin/cp", "-v", "/go/bin/terraform-provider-ldap", "/out"]
 ---> Running in 7964106bed59
Removing intermediate container 7964106bed59
 ---> 807743cb8cef
Successfully built 807743cb8cef
Successfully tagged terraform-provider-ldap:latest
================================================================
Copying the binary...
================================================================
'/go/bin/terraform-provider-ldap' -> '/out/terraform-provider-ldap'
```
