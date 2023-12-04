// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package semver

import (
	blangsemver "github.com/blang/semver/v4"

	"go.xrstf.de/rudi"
)

var (
	Functions = rudi.Functions{
		"semver": rudi.NewLiteralFunction(parseFunction, "parses a string as a semantic version").MinArgs(1).MaxArgs(1),
	}
)

func parseFunction(ctx rudi.Context, args []any) (any, error) {
	version, err := ctx.Coalesce().ToString(args[0])
	if err != nil {
		return nil, err
	}

	parsed, err := blangsemver.ParseTolerant(version)
	if err != nil {
		return nil, err
	}

	return Semver{
		Version: parsed,
	}, nil
}
