# Testing


## Prerequisites 

In order to run tests, you must have `docker` and `docker-compose` installed. For platform-specific instruction refer to the [Docker](https://hub.docker.com/search/?type=edition&offering=community) and the [Docker compose](https://docs.docker.com/compose/install/) websites. 

## Applying the Terraform recipe

In order to apply the Terraform recipe, run the following commands:

```bash
$> docker-compose up
```

to start up an OpenLDAP server and its associated management web interface;

```bash
$> terraform plan
```

to check what Terraform plans to do in order to apply the definitions as per the `main.tf` file;

```bash
$> terraform apply
```

to actually apply those definitions.

Once the Terraform commands have completem you can point your browser to http://localhost:6443 and login with username `cn=admin,dc=example,dc=com` and password `admin`; the application of the `main.tf` terraform file should have created a new `OU` called `users` and then added a `user` object named `John Doe` to it.

If you want to test what happens when an object is removed, added or modified in the `main.tf`, update the file and then run the `terraform plan` (optional!) and `terraform apply` commands again.