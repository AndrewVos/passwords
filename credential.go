package main

import (
	"fmt"
)

type Credential struct {
	SearchableContent string
	Name              string
	Site              string
	Username          string
	Password          string
}

func NewCredential(name string, site string, username string, password string) Credential {
	return Credential{
		Name:              name,
		Site:              site,
		Username:          username,
		Password:          password,
		SearchableContent: fmt.Sprintf("%v (%v)", name, site),
	}
}
