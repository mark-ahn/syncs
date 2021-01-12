package syncs_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/mark-ahn/syncs"
)

func Test_ArgMap(t *testing.T) {
	arg := syncs.ArgMap{
		"some-du": "1h30ms",
	}

	flag := arg.Parser()

	var du time.Duration
	flag.DurationVar(&du, "some-du", time.Second, "-")
	err := flag.Parse()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(du.String())
}
