package main

import (
	"crypto/tls"
	"fmt"

	"gopkg.in/ldap.v2"
)

type Config struct {
	LdapHost     string
	LdapPort     int
	UseTLS       bool
	BindUser     string
	BindPassword string
}

func (c *Config) initiateAndBind() (*ldap.Conn, error) {
	//TODO: Should we handle UDP ?
	conn, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", c.LdapHost, c.LdapPort))
	if err != nil {
		return nil, err
	}

	// Handle TLS
	if c.UseTLS {
		//TODO: Finish the TLS integration
		err = conn.StartTLS(&tls.Config{InsecureSkipVerify: true})
		if err != nil {
			return nil, err
		}
	}

	// Bind to current connection
	err = conn.Bind(c.BindUser, c.BindPassword)
	if err != nil {
		conn.Close()
		return nil, err
	}

	// Return the LDAP connection
	return conn, nil
}
