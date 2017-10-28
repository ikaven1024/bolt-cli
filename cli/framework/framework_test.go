package framework

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCommandAndArgs(t *testing.T) {
	testcases := []struct {
		name   string
		input  string
		expect string
	}{
		{
			"normal",
			"cmd opt arg",
			"cmd|opt|arg",
		},
		{
			"with space",
			`  cmd   opt   arg   `,
			"cmd|opt|arg",
		},
		{
			"with quote",
			`  "this is command" opt1 "opt 2" "opt3"`,
			"this is command|opt1|opt 2|opt3",
		},
		{
			"missing end quote",
			`  "this is command" opt1 "opt 2" "opt3 end `,
			"this is command|opt1|opt 2|opt3 end",
		},
	}

	for _, tc := range testcases {
		o := CommandAndArgs([]byte(tc.input))
		assert.Equal(t, tc.expect, string(bytes.Join(o, []byte("|"))), tc.name)
	}
}
