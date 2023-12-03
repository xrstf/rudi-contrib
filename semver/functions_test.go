// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package semver

import (
	"testing"

	blangsemver "github.com/blang/semver/v4"

	"go.xrstf.de/rudi/pkg/builtin"
	"go.xrstf.de/rudi/pkg/coalescing"
	"go.xrstf.de/rudi/pkg/testutil"
)

func TestParseFunction(t *testing.T) {
	testcases := []testutil.Testcase{
		{
			Expression: `(semver "")`,
			Invalid:    true,
		},
		{
			Expression: `(semver "foo")`,
			Invalid:    true,
		},
		{
			Expression: `(semver "1")`,
			Expected: Semver{
				Version: blangsemver.MustParse("1.0.0"),
			},
		},
		{
			Expression: `(semver "1.2.3-beta.3")`,
			Expected: Semver{
				Version: blangsemver.MustParse("1.2.3-beta.3"),
			},
		},
		{
			Expression: `(semver "v9.2")`,
			Expected: Semver{
				Version: blangsemver.MustParse("9.2.0"),
			},
		},
		{
			Expression: `(semver "v9.2").major`,
			Expected:   int64(9),
		},
		{
			Expression: `(semver "v9.2").minor`,
			Expected:   int64(2),
		},
		{
			Expression: `(semver "v9.2").patch`,
			Expected:   int64(0),
		},
		{
			Expression: `(+ (semver "v9.2").patch 1)`,
			Expected:   int64(1),
		},
		{
			Expression: `(to-string (semver "v9.2"))`,
			Expected:   "9.2.0",
		},
		{
			Expression: `(to-int (semver "v9.2"))`,
			Invalid:    true,
		},
		{
			Expression: `(eq? (semver "v9.2") "9.2.0")`,
			Expected:   true,
			Coalescer:  coalescing.NewHumane(),
		},
		{
			Expression: `(eq? (semver "v9.2") "9.2.1")`,
			Expected:   false,
			Coalescer:  coalescing.NewHumane(),
		},
		{
			Expression: `(eq? "9.2.0" (semver "v9.2"))`,
			Expected:   true,
			Coalescer:  coalescing.NewHumane(),
		},
		{
			Expression: `(eq? (semver "v9.2") "9.2.0")`,
			Invalid:    true,
		},
		{
			Expression: `(eq? (semver "v9.2") (semver "v9.2"))`,
			Expected:   true,
		},
		{
			Expression: `(eq? (semver "v9.2") (semver "v9.2.0-alpha.0"))`,
			Expected:   false,
		},
		{
			Expression: `(gt? (semver "v9.2") (semver "v9.1"))`,
			Expected:   true,
		},
		{
			Expression: `(gt? (semver "v9.2") (semver "v9.2"))`,
			Expected:   false,
		},
		{
			Expression: `(gt? (semver "v9.2") (semver "v9.2.0-alpha.0"))`,
			Expected:   true,
		},
		{
			Expression: `(gte? (semver "v9.2") (semver "v9.1"))`,
			Expected:   true,
		},
		{
			Expression: `(gte? (semver "v9.2") (semver "v9.2"))`,
			Expected:   true,
		},
		{
			Expression: `(gte? (semver "v9.2") (semver "v9.3"))`,
			Expected:   false,
		},
	}

	funcs := builtin.AllFunctions.DeepCopy().Add(Functions)

	for _, testcase := range testcases {
		testcase.Functions = funcs
		t.Run(testcase.String(), testcase.Run)
	}
}
