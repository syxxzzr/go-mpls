package main

import (
	"fmt"
	"github.com/syxxzzr/go-mpls/go-mpls"
)

func main() {
	mpls, err := go_mpls.Parse(".\\example.mpls")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%#v\n", mpls)
	}
}
