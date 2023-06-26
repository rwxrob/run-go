package run

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

// Args encapsulates the possible long-form options that can be passed
// to anything on the command line. Empty strings are valid values. To
// omit an argument ensure that it no longer exists in the map
// [delete(map,key)]. For more complex handling of options consider
// using a command-line arguments handling package.
type Args map[string]string

// List returns the expanded keys and values as a list of valid
// arguments with the keys converted into double-dash (long-form)
// arguments.
func (a Args) List() []string {
	list := []string{}
	for k, v := range a {
		k = `--` + k
		list = append(list, k, v)
	}
	return list
}

// Cmds is list of string lists to be used with the *All functions where
// each string list contains a list of arguments to be passed to Exec
// (or equivalent).
type Cmds [][]string

// Exec checks for existence of first argument as an executable on the
// system and then runs it with Go's exec.Command.Run exiting in a way that
// is supported across all architectures that Go supports. The stdin,
// stdout, and stderr are connected directly to that of the calling
// program. Sometimes this is insufficient and the UNIX-specific SysExec
// is preferred.
func Exec(args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("missing name of executable")
	}
	path, err := exec.LookPath(args[0])
	if err != nil {
		return err
	}
	cmd := exec.Command(path, args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// ExecAll simulates short-circuit logic popular in shell scripting by
// returning the first error encountered while attempting to run each
// command line from the commands list in sequential order. The index of
// the last command executed is returned as well. Note that an empty
// list of Cmds is a valid argument (where the integer returned will be
// 0).
func ExecAll(commands Cmds) (error, int) {
	var n int
	var cmd []string
	for n, cmd = range commands {
		if err := Exec(cmd...); err != nil {
			return err, n
		}
	}
	return nil, n
}

// OutErr returns the standard output of the executed command as
// a string along with any error returned.
func OutErr(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("missing name of executable")
	}
	path, err := exec.LookPath(args[0])
	if err != nil {
		return "", err
	}
	out, err := exec.Command(path, args[1:]...).Output()
	return string(out), err
}

// Out returns the standard output of the executed command as a string.
// Errors are logged after the command completes but not returned. See
// OutErr. For more current error messages during execution use Exec
// instead.
func Out(args ...string) string {
	out, err := OutErr(args...)
	// FIXME check return code and log err.Stderr if ExitError
	if err != nil {
		err, isexit := err.(*exec.ExitError)
		if isexit {
			log.Print(string(err.Stderr))
		}
		log.Print(err)
	}
	return out
}

// OutQuiet returns the standard output of the executed command as
// a string without logging any errors. Always returns a string even if
// empty. See OutErr.
func OutQuiet(args ...string) string {
	out, _ := OutErr(args...)
	return out
}

// OutAll returns the collective standard output of the executed
// commands as a string or the first error encountered along with the
// index in Cmds of the command arguments that produced the error. See
// OutErr.
func OutAll(commands Cmds) (string, error, int) {
	var n int
	var cmd []string
	var buf string
	for n, cmd = range commands {
		out, err := OutErr(cmd...)
		buf += out
		if err != nil {
			return buf, err, n
		}
	}
	return buf, nil, n
}
