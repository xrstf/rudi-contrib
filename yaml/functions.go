// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package yaml

import (
	"fmt"

	"go.xrstf.de/rudi/pkg/eval"
	"go.xrstf.de/rudi/pkg/eval/types"
	"go.xrstf.de/rudi/pkg/lang/ast"

	yamlv3 "gopkg.in/yaml.v3"
)

var (
	Functions = types.Functions{
		"to-yaml":   types.BasicFunction(toYamlFunction, "encodes the given value as YAML"),
		"from-yaml": types.BasicFunction(fromYamlFunction, "decodes a YAML string into a Go value"),
	}
)

func toYamlFunction(ctx types.Context, args []ast.Expression) (any, error) {
	if size := len(args); size != 1 {
		return nil, fmt.Errorf("expected 1 argument, got %d", size)
	}

	_, data, err := eval.EvalExpression(ctx, args[0])
	if err != nil {
		return nil, err
	}

	encoded, err := yamlv3.Marshal(data)
	if err != nil {
		return nil, err
	}

	return string(encoded), nil
}

func fromYamlFunction(ctx types.Context, args []ast.Expression) (any, error) {
	if size := len(args); size != 1 {
		return nil, fmt.Errorf("expected 1 argument, got %d", size)
	}

	_, data, err := eval.EvalExpression(ctx, args[0])
	if err != nil {
		return nil, err
	}

	dataString, err := ctx.Coalesce().ToString(data)
	if err != nil {
		return nil, err
	}

	var result any
	if err := yamlv3.Unmarshal([]byte(dataString), &result); err != nil {
		return nil, err
	}

	return result, nil
}
