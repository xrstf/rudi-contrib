// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package set

import (
	"errors"
	"fmt"

	"k8s.io/apimachinery/pkg/util/sets"

	"go.xrstf.de/rudi"
	"go.xrstf.de/rudi/pkg/runtime/types"
)

var (
	Functions = rudi.Functions{
		"new-set":     rudi.NewFunctionBuilder(newEmptySetFunction, newSetFunction).WithDescription("create a set filled with the given values").Build(),
		"new-key-set": rudi.NewFunctionBuilder(keySetFunction).WithDescription("create a set filled with the keys of an object").Build(),

		"set-delete":       rudi.NewFunctionBuilder(setDeleteFunction).WithDescription("returns a copy of the set with the given values removed from it").Build(),
		"set-diff":         rudi.NewFunctionBuilder(setDifferenceFunction).WithDescription("returns the difference between two sets").Build(),
		"set-insert":       rudi.NewFunctionBuilder(setInsertFunction).WithDescription("returns a copy of the set with the newly added values inserted to it").Build(),
		"set-intersection": rudi.NewFunctionBuilder(setIntersectionFunction).WithDescription("returns the insersection of two sets").Build(),
		"set-size":         rudi.NewFunctionBuilder(setLenFunction).WithDescription("returns the number of values in the set").Build(),
		"set-symdiff":      rudi.NewFunctionBuilder(setSymmetricDifferenceFunction).WithDescription("returns the symmetric difference between two sets").Build(),
		"set-union":        rudi.NewFunctionBuilder(setUnionFunction).WithDescription("returns the union of two or more sets").Build(),

		"set-eq?":          rudi.NewFunctionBuilder(setEqualFunction).WithDescription("returns true if two sets hold the same values").Build(),
		"set-has?":         rudi.NewFunctionBuilder(setHasFunction).WithDescription("returns true if the set contains _all_ of the given values").Build(),
		"set-has-any?":     rudi.NewFunctionBuilder(setHasAnyFunction).WithDescription("returns true if the set contains _any_ of the given values").Build(),
		"set-superset-of?": rudi.NewFunctionBuilder(setIsSupersetFunction).WithDescription("returns true if the other set is a superset of the base set").Build(),
	}
)

func newEmptySetFunction() (any, error) {
	return sets.New[string](), nil
}

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

func setDeleteFunction(ctx types.Context, target any, itemsToRemove ...any) (any, error) {
	s, ok := target.(sets.Set[string])
	if !ok {
		return nil, fmt.Errorf("argument #0: not a set, but %T", target)
	}

	strs, err := toStrings(ctx, itemsToRemove...)
	if err != nil {
		return nil, err
	}

	// NB: Remove from a clone of the set; removing inplace happens via bang modifier magic
	// (i.e. "(set-delete! $myset 1 2 3)")
	return s.Clone().Delete(strs...), nil
}

func setLenFunction(ctx types.Context, target any) (any, error) {
	s, ok := target.(sets.Set[string])
	if !ok {
		return nil, fmt.Errorf("argument #0: not a set, but %T", target)
	}

	return s.Len(), nil
}

func setHasFunction(ctx types.Context, target any, items ...any) (any, error) {
	s, ok := target.(sets.Set[string])
	if !ok {
		return nil, fmt.Errorf("argument #0: not a set, but %T", target)
	}

	strs, err := toStrings(ctx, items...)
	if err != nil {
		return nil, err
	}

	return s.HasAll(strs...), nil
}

func setHasAnyFunction(ctx types.Context, target any, items ...any) (any, error) {
	s, ok := target.(sets.Set[string])
	if !ok {
		return nil, fmt.Errorf("argument #0: not a set, but %T", target)
	}

	strs, err := toStrings(ctx, items...)
	if err != nil {
		return nil, err
	}

	return s.HasAny(strs...), nil
}

type setsFunc func(a, b sets.Set[string]) (any, error)

func callFuncOnSets(a any, b any, f setsFunc) (any, error) {
	setA, ok := a.(sets.Set[string])
	if !ok {
		return nil, fmt.Errorf("argument #0: not a set, but %T", a)
	}

	setB, ok := b.(sets.Set[string])
	if !ok {
		return nil, fmt.Errorf("argument #1: not a set, but %T", b)
	}

	return f(setA, setB)
}

func setEqualFunction(target any, other any) (any, error) {
	return callFuncOnSets(target, other, func(a, b sets.Set[string]) (any, error) {
		return a.Equal(b), nil
	})
}

func setIntersectionFunction(ctx types.Context, target any, other any) (any, error) {
	return callFuncOnSets(target, other, func(a, b sets.Set[string]) (any, error) {
		return a.Intersection(b), nil
	})
}

func setDifferenceFunction(ctx types.Context, target any, other any) (any, error) {
	return callFuncOnSets(target, other, func(a, b sets.Set[string]) (any, error) {
		return a.Difference(b), nil
	})
}

func setSymmetricDifferenceFunction(ctx types.Context, target any, other any) (any, error) {
	return callFuncOnSets(target, other, func(a, b sets.Set[string]) (any, error) {
		return a.SymmetricDifference(b), nil
	})
}

func setIsSupersetFunction(ctx types.Context, target any, other any) (any, error) {
	return callFuncOnSets(target, other, func(a, b sets.Set[string]) (any, error) {
		return a.IsSuperset(b), nil
	})
}

func setUnionFunction(ctx types.Context, target any, others ...any) (any, error) {
	acc, ok := target.(sets.Set[string])
	if !ok {
		return nil, fmt.Errorf("argument #0: not a set, but %T", target)
	}

	for i, otherSet := range others {
		toUnionize, ok := otherSet.(sets.Set[string])
		if !ok {
			return nil, fmt.Errorf("argument #%d: not a set, but %T", i+1, otherSet)
		}

		acc = acc.Union(toUnionize)
	}

	return acc, nil
}

func insertMany(ctx types.Context, s sets.Set[string], vals ...any) (any, error) {
	strs, err := toStrings(ctx, vals...)
	if err != nil {
		return nil, err
	}

	s.Insert(strs...)

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
