// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package yaml

import (
	yamlv3 "gopkg.in/yaml.v3"

	"go.xrstf.de/rudi"
)

var (
	Functions = rudi.Functions{
		"to-yaml":   rudi.NewFunctionBuilder(toYamlFunction).WithDescription("encodes the given value as YAML").Build(),
		"from-yaml": rudi.NewFunctionBuilder(fromYamlFunction).WithDescription("decodes a YAML string into a Go value").Build(),
	}
)

func toYamlFunction(val any) (any, error) {
	encoded, err := yamlv3.Marshal(val)
	if err != nil {
		return nil, err
	}

	return string(encoded), nil
}

func fromYamlFunction(encoded string) (any, error) {
	var result any
	if err := yamlv3.Unmarshal([]byte(encoded), &result); err != nil {
		return nil, err
	}

	return result, nil
}
