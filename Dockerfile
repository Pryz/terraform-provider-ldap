FROM golang:1.11.3-stretch
COPY . /go/src/terraform-provider-ldap
WORKDIR /go/src/terraform-provider-ldap
RUN go get .
RUN CGO_ENABLED=0 GOOS=linux go install -a -ldflags '-extldflags "-static"' .
ENTRYPOINT ["/bin/cp", "-v", "/go/bin/terraform-provider-ldap", "/out"]
