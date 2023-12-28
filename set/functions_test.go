// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package set

import (
	"testing"

	"k8s.io/apimachinery/pkg/util/sets"

	"go.xrstf.de/rudi/pkg/builtin"
	"go.xrstf.de/rudi/pkg/testutil"
)

func TestNewSetFunction(t *testing.T) {
	testcases := []testutil.Testcase{
		{
			Expression: `(new-set true)`,
			Invalid:    true, // because of strict coalescing
		},
		{
			Expression: `(new-set {a "b"})`,
			Invalid:    true,
		},
		{
			Expression: `(new-set)`,
			Expected:   sets.New[string](),
		},
		{
			Expression: `(new-set "a" "b")`,
			Expected:   sets.New[string]("a", "b"),
		},
		{
			Expression: `(new-set ["a" "b"])`,
			Expected:   sets.New[string]("a", "b"),
		},
		{
			Expression: `(new-set "a" [])`,
			Expected:   sets.New[string]("a"),
		},
		{
			Expression: `(new-set "a" "")`,
			Expected:   sets.New[string]("a", ""),
		},
		{
			Expression: `(new-set "a" ["b" "c" ""])`,
			Expected:   sets.New[string]("a", "b", "c", ""),
		},
		{
			Expression: `(new-set "a" [["b"]])`,
			Invalid:    true, // do not recurse more than one level
		},
		{
			// do not explode if a value occurs multiple times
			Expression: `(new-set ["a" "b" "a"])`,
			Expected:   sets.New[string]("a", "b"),
		},
	}

	funcs := builtin.SafeFunctions.DeepCopy().Add(Functions)

	for _, testcase := range testcases {
		testcase.Functions = funcs
		t.Run(testcase.String(), testcase.Run)
	}
}

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

	funcs := builtin.SafeFunctions.DeepCopy().Add(Functions)

	for _, testcase := range testcases {
		testcase.Functions = funcs
		t.Run(testcase.String(), testcase.Run)
	}
}

func TestSetInsertFunction(t *testing.T) {
	testcases := []testutil.Testcase{
		{
			Expression: `(set-insert (new-set "a" "b") "c")`,
			Expected:   sets.New[string]("a", "b", "c"),
		},
		{
			Expression: `(set-insert (new-set "a" "b") ["c" ""] "d" "a")`,
			Expected:   sets.New[string]("a", "b", "c", "", "d"),
		},

		// do not modify in-place

		{
			Expression: `(set! $s (new-set "a" "b")) (set-insert $s "c")`,
			Expected:   sets.New[string]("a", "b", "c"),
		},
		{
			Expression: `(set! $s (new-set "a" "b")) (set-insert $s "c") $s`,
			Expected:   sets.New[string]("a", "b"),
		},

		// modify in-place

		{
			Expression: `(set! $s (new-set "a" "b")) (set-insert! $s "c")`,
			Expected:   sets.New[string]("a", "b", "c"),
		},
		{
			Expression: `(set! $s (new-set "a" "b")) (set-insert! $s "c") $s`,
			Expected:   sets.New[string]("a", "b", "c"),
		},
	}

	funcs := builtin.SafeFunctions.DeepCopy().Add(Functions)

	for _, testcase := range testcases {
		testcase.Functions = funcs
		t.Run(testcase.String(), testcase.Run)
	}
}

func TestSetDeleteFunction(t *testing.T) {
	testcases := []testutil.Testcase{
		{
			Expression: `(set-delete (new-set "a" "b") "c")`,
			Expected:   sets.New[string]("a", "b"),
		},
		{
			Expression: `(set-delete (new-set "a" "b") "b")`,
			Expected:   sets.New[string]("a"),
		},
		{
			Expression: `(set-delete (new-set "a" "b" "") ["c" ""] "d" "a")`,
			Expected:   sets.New[string]("b"),
		},

		// do not modify in-place

		{
			Expression: `(set! $s (new-set "a" "b")) (set-delete $s "b")`,
			Expected:   sets.New[string]("a"),
		},
		{
			Expression: `(set! $s (new-set "a" "b")) (set-delete $s "b") $s`,
			Expected:   sets.New[string]("a", "b"),
		},

		// modify in-place

		{
			Expression: `(set! $s (new-set "a" "b")) (set-delete! $s "b")`,
			Expected:   sets.New[string]("a"),
		},
		{
			Expression: `(set! $s (new-set "a" "b")) (set-delete! $s "b") $s`,
			Expected:   sets.New[string]("a"),
		},
	}

	funcs := builtin.SafeFunctions.DeepCopy().Add(Functions)

	for _, testcase := range testcases {
		testcase.Functions = funcs
		t.Run(testcase.String(), testcase.Run)
	}
}

func TestSetSizeFunction(t *testing.T) {
	testcases := []testutil.Testcase{
		{
			Expression: `(set-size "nope")`,
			Invalid:    true,
		},
		{
			Expression: `(set-size (new-set "a" "b"))`,
			Expected:   2,
		},
		{
			Expression: `(set-size (set-delete (new-set "a") "a"))`,
			Expected:   0,
		},
	}

	funcs := builtin.SafeFunctions.DeepCopy().Add(Functions)

	for _, testcase := range testcases {
		testcase.Functions = funcs
		t.Run(testcase.String(), testcase.Run)
	}
}

func TestSetHasFunction(t *testing.T) {
	testcases := []testutil.Testcase{
		{
			Expression: `(set-has? "nope")`,
			Invalid:    true,
		},
		{
			Expression: `(set-has? (new-set "a" "b"))`,
			Invalid:    true,
		},
		{
			Expression: `(set-has? (new-set "a" "b") "c")`,
			Expected:   false,
		},
		{
			Expression: `(set-has? (new-set "a" "b") "a")`,
			Expected:   true,
		},
		{
			Expression: `(set-has? (new-set "a" "b") ["a"])`,
			Expected:   true,
		},
		{
			Expression: `(set-has? (new-set "a" "b") "a" "b")`,
			Expected:   true,
		},
		{
			Expression: `(set-has? (new-set "a" "b") "a" "c")`,
			Expected:   false,
		},
	}

	funcs := builtin.SafeFunctions.DeepCopy().Add(Functions)

	for _, testcase := range testcases {
		testcase.Functions = funcs
		t.Run(testcase.String(), testcase.Run)
	}
}

func TestSetHasAnyFunction(t *testing.T) {
	testcases := []testutil.Testcase{
		{
			Expression: `(set-has-any? "nope")`,
			Invalid:    true,
		},
		{
			Expression: `(set-has-any? (new-set "a" "b"))`,
			Invalid:    true,
		},
		{
			Expression: `(set-has-any? (new-set "a" "b") "c")`,
			Expected:   false,
		},
		{
			Expression: `(set-has-any? (new-set "a" "b") "a")`,
			Expected:   true,
		},
		{
			Expression: `(set-has-any? (new-set "a" "b") ["a"])`,
			Expected:   true,
		},
		{
			Expression: `(set-has-any? (new-set "a" "b") "a" "b")`,
			Expected:   true,
		},
		{
			Expression: `(set-has-any? (new-set "a" "b") "a" "c")`,
			Expected:   true,
		},
	}

	funcs := builtin.SafeFunctions.DeepCopy().Add(Functions)

	for _, testcase := range testcases {
		testcase.Functions = funcs
		t.Run(testcase.String(), testcase.Run)
	}
}

func TestSetListFunction(t *testing.T) {
	testcases := []testutil.Testcase{
		{
			Expression: `(set-list "nope")`,
			Invalid:    true,
		},
		{
			Expression: `(set-list (new-set "b" "a"))`,
			Expected:   []any{"a", "b"},
		},
	}

	funcs := builtin.SafeFunctions.DeepCopy().Add(Functions)

	for _, testcase := range testcases {
		testcase.Functions = funcs
		t.Run(testcase.String(), testcase.Run)
	}
}

func TestSetUnionFunction(t *testing.T) {
	testcases := []testutil.Testcase{
		{
			Expression: `(set-union "nope")`,
			Invalid:    true,
		},
		{
			Expression: `(set-union (new-set "a" "b"))`,
			Invalid:    true,
		},
		{
			Expression: `(set-union (new-set "a" "b") "nope")`,
			Invalid:    true,
		},
		{
			Expression: `(set-union "nope" (new-set "a" "b"))`,
			Invalid:    true,
		},
		{
			Expression: `(set-union (new-set "a" "b") (new-set "c"))`,
			Expected:   sets.New[string]("a", "b", "c"),
		},
		{
			Expression: `(set-union (new-set "a" "b") (new-set "a"))`,
			Expected:   sets.New[string]("a", "b"),
		},
		{
			Expression: `(set-union (new-set "a" "b") (new-set "a") (new-set "d"))`,
			Expected:   sets.New[string]("a", "b", "d"),
		},
	}

	funcs := builtin.SafeFunctions.DeepCopy().Add(Functions)

	for _, testcase := range testcases {
		testcase.Functions = funcs
		t.Run(testcase.String(), testcase.Run)
	}
}
