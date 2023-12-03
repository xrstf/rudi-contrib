// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package semver

import (
	"fmt"

	"go.xrstf.de/rudi/pkg/eval"
	"go.xrstf.de/rudi/pkg/eval/types"
	"go.xrstf.de/rudi/pkg/lang/ast"

	blangsemver "github.com/blang/semver/v4"
)

var (
	Functions = types.Functions{
		"semver": types.BasicFunction(parseFunction, "parses a string as a semantic version"),
	}
)

func parseFunction(ctx types.Context, args []ast.Expression) (any, error) {
	if size := len(args); size != 1 {
		return nil, fmt.Errorf("expected 1 argument, got %d", size)
	}

	_, version, err := eval.EvalExpression(ctx, args[0])
	if err != nil {
		return nil, err
	}

	versionString, err := ctx.Coalesce().ToString(version)
	if err != nil {
		return nil, err
	}

	parsed, err := blangsemver.ParseTolerant(versionString)
	if err != nil {
		return nil, err
	}

	return Semver{
		Version: parsed,
	}, nil
}
