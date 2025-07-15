package go_mpls

import (
	"fmt"
	"os"
	"testing"
)

// NOTICE: the environment variable named `MPLS_PATH` which pointed to the *.mpls file should be set before test
func TestParse(t *testing.T) {
	mpls, err := Parse(os.Getenv("MPLS_PATH"))
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%#v\n", mpls)
	}
}
