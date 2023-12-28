// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package set

import (
	"errors"
	"fmt"

	"k8s.io/apimachinery/pkg/util/sets"

	"go.xrstf.de/rudi"
	"go.xrstf.de/rudi/pkg/eval/types"
)

var (
	Functions = rudi.Functions{
		"new-set":     rudi.NewFunctionBuilder(newSetFunction).WithDescription("create a set filled with the given values").Build(),
		"new-key-set": rudi.NewFunctionBuilder(keySetFunction).WithDescription("create a set filled with the keys of an object").Build(),
		"set-insert":  rudi.NewFunctionBuilder(setInsertFunction).WithDescription("returns a copy of the set with the newly added values inserted to it").Build(),
		"set-delete":  rudi.NewFunctionBuilder(setDeleteFunction).WithDescription("returns a copy of the set with the given values removed from it").Build(),
	}
)

func newSetFunction(ctx types.Context, vals ...any) (any, error) {
	return insertMany(ctx, sets.New[string](), vals...)
}

func keySetFunction(val map[string]any) (any, error) {
	return sets.KeySet[string](val), nil
}

func setInsertFunction(ctx types.Context, target any, newItems ...any) (any, error) {
	s, ok := target.(sets.Set[string])
	if !ok {
		return nil, fmt.Errorf("argument #0: not a set, but %T", target)
	}

	// NB: Insert into a clone of the set; adding inplace happens via bang modifier magic
	// (i.e. "(set-insert! $myset 1 2 3)")
	return insertMany(ctx, s.Clone(), newItems...)
}

func setDeleteFunction(ctx types.Context, target any, ttemsToRemove ...any) (any, error) {
	s, ok := target.(sets.Set[string])
	if !ok {
		return nil, fmt.Errorf("argument #0: not a set, but %T", target)
	}

	// NB: Remove from a clone of the set; removing inplace happens via bang modifier magic
	// (i.e. "(set-delete! $myset 1 2 3)")
	return deleteMany(ctx, s.Clone(), ttemsToRemove...)
}

func insertMany(ctx types.Context, s sets.Set[string], vals ...any) (any, error) {
	strs, err := toStrings(ctx, vals...)
	if err != nil {
		return nil, err
	}

	s.Insert(strs...)

	return s, nil
}

func deleteMany(ctx types.Context, s sets.Set[string], vals ...any) (any, error) {
	strs, err := toStrings(ctx, vals...)
	if err != nil {
		return nil, err
	}

	s.Delete(strs...)

	return s, nil
}

func toStrings(ctx types.Context, vals ...any) ([]string, error) {
	result := []string{}

	for _, v := range vals {
		// This is purposefully not recursive so we do not run into unexpected situations.
		vec, err := ctx.Coalesce().ToVector(v)
		if err == nil {
			for _, v := range vec {
				str, err := ctx.Coalesce().ToString(v)
				if err != nil {
					return nil, errors.New("argument vector contains non-string")
				}
				result = append(result, str)
			}

			continue
		}

		str, err := ctx.Coalesce().ToString(v)
		if err != nil {
			return nil, fmt.Errorf("argument is neither vector nor string, but %T", v)
		}

		result = append(result, str)
	}

	return result, nil
}
