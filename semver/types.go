// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package semver

import (
	"fmt"
	"strings"

	"go.xrstf.de/rudi/pkg/coalescing"
	"go.xrstf.de/rudi/pkg/deepcopy"
	"go.xrstf.de/rudi/pkg/equality"
	"go.xrstf.de/rudi/pkg/pathexpr"

	blangsemver "github.com/blang/semver/v4"
)

type Semver struct {
	Version blangsemver.Version
}

var (
	_ pathexpr.ObjectReader            = Semver{}
	_ pathexpr.ObjectWriter            = Semver{}
	_ deepcopy.Copier                  = Semver{}
	_ coalescing.CustomStringCoalescer = Semver{}
	_ equality.Comparer                = Semver{}
)

// GetObjectKey implements pathexpr.ObjectWriter.
func (v Semver) GetObjectKey(name string) (any, error) {
	switch strings.ToLower(name) {
	case "major":
		return int64(v.Version.Major), nil
	case "minor":
		return int64(v.Version.Minor), nil
	case "patch":
		return int64(v.Version.Patch), nil
	default:
		return nil, fmt.Errorf("unknown property %q", name)
	}
}

func toInteger(val any) (int64, bool) {
	switch asserted := val.(type) {
	case int:
		return int64(asserted), true
	case int32:
		return int64(asserted), true
	case int64:
		return asserted, true
	case uint64:
		return int64(asserted), true
	default:
		return 0, false
	}
}

// SetObjectKey implements pathexpr.ObjectWriter.
func (v Semver) SetObjectKey(name string, value any) (any, error) {
	switch strings.ToLower(name) {
	case "major":
		val, ok := toInteger(value)
		if !ok {
			return nil, fmt.Errorf("cannot set major version to %T", value)
		}
		v.Version.Major = uint64(val)
	case "minor":
		val, ok := toInteger(value)
		if !ok {
			return nil, fmt.Errorf("cannot set minor version to %T", value)
		}
		v.Version.Minor = uint64(val)
	case "patch":
		val, ok := toInteger(value)
		if !ok {
			return nil, fmt.Errorf("cannot set patch version to %T", value)
		}
		v.Version.Patch = uint64(val)
	default:
		return nil, fmt.Errorf("unknown property %q", name)
	}

	return v, nil
}

// DeepCopy implements deepcopy.Copier.
func (v Semver) DeepCopy() (any, error) {
	pres := make([]blangsemver.PRVersion, len(v.Version.Pre))
	copy(pres, v.Version.Pre)

	builds := make([]string, len(v.Version.Build))
	copy(builds, v.Version.Build)

	return Semver{
		Version: blangsemver.Version{
			Major: v.Version.Major,
			Minor: v.Version.Minor,
			Patch: v.Version.Patch,
			Pre:   pres,
			Build: builds,
		},
	}, nil
}

// CoalesceToString implements coalescing.CustomStringCoalescer.
func (v Semver) CoalesceToString(_ coalescing.Coalescer) (string, error) {
	return v.Version.String(), nil
}

// Equal implements equality.Comparer.
func (v Semver) Compare(other any) (int, error) {
	otherV, ok := other.(Semver)
	if !ok {
		return 0, equality.ErrIncompatibleTypes
	}

	return v.Version.Compare(otherV.Version), nil
}
