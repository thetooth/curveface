package main

import (
	"testing"

	"github.com/thetooth/curveface"
)

func Test(t *testing.T) {
	if res, err := curveface.GetGface(
		"127.0.0.1:9001",
		"p4dIvqYyJSn1K9OUpxcB/ouVrX8KmT0CVwc/L/hPMHWo=", "PiUVHEeVQvYVowItse88kConVCOBMiLPlsd5y15Trb2c=",
		"P9/NIhsJvQAvnxoc2177O/3aIzUHYOjcack2wkkpTn2Q=",
	); err != nil {
		t.Fatal(err)
	} else {
		t.Logf("%s", res)
	}
}
