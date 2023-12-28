// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package set

import (
	"testing"

	"k8s.io/apimachinery/pkg/util/sets"

	"go.xrstf.de/rudi/pkg/builtin"
	"go.xrstf.de/rudi/pkg/testutil"
)

func TestNewKeySetFunction(t *testing.T) {
	testcases := []testutil.Testcase{
		{
			Expression: `(new-key-set "")`,
			Invalid:    true,
		},
		{
			Expression: `(new-key-set [1 2 3])`,
			Invalid:    true,
		},
		{
			Expression: `(new-key-set true)`,
			Invalid:    true,
		},
		{
			Expression: `(new-key-set {})`,
			Expected:   sets.New[string](),
		},
		{
			Expression: `(new-key-set {a "b" c "d"})`,
			Expected:   sets.New[string]("a", "c"),
		},
		{
			// do not explode if an object literal contains the same key twice
			Expression: `(new-key-set {a "b" c "d" c "x"})`,
			Expected:   sets.New[string]("a", "c"),
		},
	}

	funcs := builtin.Functions.DeepCopy().Add(Functions)

	for _, testcase := range testcases {
		testcase.Functions = funcs
		t.Run(testcase.String(), testcase.Run)
	}
}
