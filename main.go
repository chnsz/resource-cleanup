package main

import (
	rc "github.com/chnsz/resource-cleanup/pkg"
)

func main() {
	cleaners := []rc.ResourceQuery{
		&rc.VpcSubnet{},
	}

	for _, c := range cleaners {
		name := c.GetName()
		ids := c.QueryIds()
		rc.Clean(name, ids)
	}
}
