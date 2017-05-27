# Testing

```
sudo docker-compose up
terraform plan
```
Then point your browser to http://localhost:6443 and login with username
"cn=admin,dc=example,dc=com" and password "admin"; the application of the
main.tf terraform file creates an OU called "users", then create a user under
that OU.