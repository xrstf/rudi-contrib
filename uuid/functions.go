// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package uuid

import (
	guuid "github.com/google/uuid"

	"go.xrstf.de/rudi"
)

var (
	Functions = rudi.Functions{
		"uuidv4": rudi.NewFunctionBuilder(newUUIDv4Function).WithDescription("returns a new, randomly generated v4 UUID").Build(),
	}
)

func newUUIDv4Function() (any, error) {
	id, err := guuid.NewRandom()
	if err != nil {
		return nil, err
	}

	return id.String(), nil
}
