// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package semver

import (
	blangsemver "github.com/blang/semver/v4"

	"go.xrstf.de/rudi"
)

var (
	Functions = rudi.Functions{
		"semver": rudi.NewFunctionBuilder(parseFunction).WithDescription("parses a string as a semantic version").Build(),
	}
)

func parseFunction(version string) (any, error) {
	parsed, err := blangsemver.ParseTolerant(version)
	if err != nil {
		return nil, err
	}

	return Semver{
		Version: parsed,
	}, nil
}
