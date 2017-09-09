package opentis100

import (
	"fmt"
	"testing"
)

type testCase struct {
	op         string
	arg1       string
	arg2       string
	expectedOp operation
}

func TestNewInstruction(t *testing.T) {

	cases := []testCase{
		testCase{op: "mov", arg1: "1", arg2: "acc", expectedOp: mov},
		testCase{op: "sub", arg1: "1", expectedOp: sub},
		testCase{op: "add", arg1: "1", expectedOp: add},
	}

	for _, cas := range cases {
		t.Run(fmt.Sprintf("New %s instruction", cas.op), func(t *testing.T) {
			i := newInstruction(cas.op, cas.arg1, cas.arg2)

			if i.Operation != cas.expectedOp {
				t.Errorf("Expected instruction %s, got %s instead", cas.expectedOp, i.Operation)
			}
		})
	}
}
