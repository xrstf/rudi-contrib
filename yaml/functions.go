// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package yaml

import (
	yamlv3 "gopkg.in/yaml.v3"

	"go.xrstf.de/rudi"
)

var (
	Functions = rudi.Functions{
		"to-yaml":   rudi.NewLiteralFunction(toYamlFunction, "encodes the given value as YAML").MinArgs(1).MaxArgs(1),
		"from-yaml": rudi.NewLiteralFunction(fromYamlFunction, "decodes a YAML string into a Go value").MinArgs(1).MaxArgs(1),
	}
)

func toYamlFunction(ctx rudi.Context, args []any) (any, error) {
	encoded, err := yamlv3.Marshal(args[0])
	if err != nil {
		return nil, err
	}

	return string(encoded), nil
}

func fromYamlFunction(ctx rudi.Context, args []any) (any, error) {
	data, err := ctx.Coalesce().ToString(args[0])
	if err != nil {
		return nil, err
	}

	var result any
	if err := yamlv3.Unmarshal([]byte(data), &result); err != nil {
		return nil, err
	}

	return result, nil
}
