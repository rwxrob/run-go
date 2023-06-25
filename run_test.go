package run_test

import (
	"fmt"

	"github.com/rwxrob/run-go"
)

func ExampleExecAll() {

	err, n := run.ExecAll(run.Cmds{{`ls`, `-d`, `/tmp`}, {`true`}, {`echo`, `wow`}})
	fmt.Println(err, n)

	// Output:
	// /tmp
	// wow
	// <nil> 2
}

func ExampleExecAll_mid_Fail() {

	err, n := run.ExecAll(run.Cmds{{`ls`, `-d`, `/tmp`}, {`false`}, {`echo`, `it works`}})
	fmt.Println(err, n)

	// Output:
	// /tmp
	// exit status 1 1
}

func ExampleExecAll_first_Fail() {

	err, n := run.ExecAll(run.Cmds{{`ls`, `notme`}, {`true`}, {`echo`, `it works`}})
	fmt.Println(err, n)

	// Output:
	// exit status 2 0
}

func ExampleExecAll_last_Fail() {

	err, n := run.ExecAll(run.Cmds{{`ls`, `-d`, `/tmp`}, {`true`}, {`ls`, `bork`}})
	fmt.Println(err, n)

	// Output:
	// /tmp
	// exit status 2 2
}

func ExampleOutErr() {
	out, err := run.OutErr(`ls`, `-d`, `/tmp`)
	fmt.Printf("%q %v", out, err)
	// Output:
	// "/tmp\n" <nil>
}

func ExampleOutErr_with_Error() {
	out, err := run.OutErr(`ls`, `-d`, `/nopenothear`)
	fmt.Printf("%q %v", out, err)
	// Output:
	// "" exit status 2
}

func ExampleOut() {
	out := run.Out(`ls`, `-d`, `/tmp`)
	fmt.Printf("%q", out)
	// Output:
	// "/tmp\n"
}

func ExampleOut_with_Error() {
	out := run.Out(`ls`, `-d`, `/nopenothear`)
	// note the output to stderr
	fmt.Printf("%q", out)
	// Output:
	// ""
}

func ExampleOutQuiet() {
	out := run.OutQuiet(`ls`, `-d`, `/tmp`)
	fmt.Printf("%q", out)
	// Output:
	// "/tmp\n"
}

func ExampleOutQuiet_with_Error() {
	out := run.OutQuiet(`ls`, `-d`, `/nopenothear`)
	fmt.Printf("%q", out)
	// Output:
	// ""
}

func ExampleOutAll() {

	buf, err, n := run.OutAll(run.Cmds{{`ls`, `-d`, `/tmp`}, {`true`}, {`echo`, `wow`}})
	fmt.Printf("%q %v %v", buf, err, n)

	// Output:
	// "/tmp\nwow\n" <nil> 2
}
