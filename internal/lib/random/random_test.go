package random

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestNewRandomString(t *testing.T) {
    tests := []struct {
        name    string 	// test name
        arg     int 	// argument 
        want    int 	// expected value
    }{
        // testcases
		{
			"length 0",
			0,
			0,
		},
		{
			"length 1",
			1,
			1,
		},
		{
			"length 5",
			5,
			5,
		},
		{
			"length 10",
			10,
			10,
		},
    }

    // calling testing function for each test case  
    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, len(NewRandomString(test.arg)), test.want)
        })
    }
}
