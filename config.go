package main

import (
	"crypto/tls"
	"fmt"

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

func (c *Config) initiateAndBind() (*ldap.Conn, error) {
	// TODO: should we handle UDP ?
	connection, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", c.LDAPHost, c.LDAPPort))
	if err != nil {
		return nil, err
	}

	// handle TLS
	if c.UseTLS {
		//TODO: Finish the TLS integration
		err = connection.StartTLS(&tls.Config{InsecureSkipVerify: true})
		if err != nil {
			return nil, err
		}
	}

	// bind to current connection
	err = connection.Bind(c.BindUser, c.BindPassword)
	if err != nil {
		connection.Close()
		return nil, err
	}

	// return the LDAP connection
	return connection, nil
}
