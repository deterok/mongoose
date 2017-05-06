package mongoose

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubCommandParse(t *testing.T) {
	root := Command{MinArgs: 2}
	root.Flags().BoolP("flag", "", false, "A simple bollean flag.")

	subCmd := &Command{
		Name:    "Sub-Command",
		MinArgs: 1,
	}

	subCmd.Flags().StringP("sub-flag", "", "", "A simple string flag.")

	root.AddCommand(subCmd)

	subSubCmd := &Command{
		Name: "Sub-Sub-Command",
	}
	root.AddCommand(subSubCmd)

	args := []string{"my_app", "--flag", "arg1", "arg2",
		"sub-command", "--sub-flag", "sub_flag_value", "sub_arg1",
		"sub-sub-command"}
	root.Execute(args[1:])

	assert.Equal(t, []string{"arg1", "arg2"}, root.Args)
	assert.Equal(t, []string{"sub-command", "--sub-flag", "sub_flag_value",
		"sub_arg1", "sub-sub-command"}, root.Tail)

	assert.Equal(t, []string{"sub_arg1"}, subCmd.Args)
	assert.Equal(t, []string{"sub-sub-command"}, subCmd.Tail)
}
