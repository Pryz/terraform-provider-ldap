BINARY=terraform-provider-ldap
TEST_ENV := LDAP_HOST=localhost LDAP_PORT=389 LDAP_BIND_USER="cn=admin,dc=example,dc=com" LDAP_BIND_PASSWORD=admin

.DEFAULT_GOAL: $(BINARY)

$(BINARY):
	go build -o bin/$(BINARY)

test:
	go test -v

docker_test:
	$(TEST_ENV) go test -v
