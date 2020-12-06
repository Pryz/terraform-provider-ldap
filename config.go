package main

import (
	"crypto/tls"
	"fmt"
	"sync"

	"gopkg.in/ldap.v2"
)

// Config is the set of parameters needed to configure the LDAP provider.
type Config struct {
	LDAPHost     string
	LDAPPort     int
	UseTLS       bool
	BindUser     string
	BindPassword string
}

// The lazily initialized response from the first initiateAndBind attempt.
var initiateAndBindResponse *InitiateAndBindResponse
var once sync.Once

type InitiateAndBindResponse struct {
	Connection *ldap.Conn
	Err        error
}

func (c *Config) initiateAndBind() (*ldap.Conn, error) {
	once.Do(func() {
		initiateAndBindResponse = &InitiateAndBindResponse{}

		// TODO: should we handle UDP ?
		connection, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", c.LDAPHost, c.LDAPPort))
		if err != nil {
			initiateAndBindResponse.Err = err
			return
		}

		// handle TLS
		if c.UseTLS {
			//TODO: Finish the TLS integration
			err = connection.StartTLS(&tls.Config{InsecureSkipVerify: true})
			if err != nil {
				connection.Close()
				initiateAndBindResponse.Err = err
				return
			}
		}

		// bind to current connection
		err = connection.Bind(c.BindUser, c.BindPassword)
		if err != nil {
			connection.Close()
			initiateAndBindResponse.Err = err
		}

		// return the LDAP connection
		initiateAndBindResponse.Connection = connection
	})
	return initiateAndBindResponse.Connection, initiateAndBindResponse.Err
}
