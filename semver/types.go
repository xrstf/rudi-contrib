// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package semver

import (
	blangsemver "github.com/blang/semver/v4"

	"go.xrstf.de/rudi/pkg/coalescing"
	"go.xrstf.de/rudi/pkg/deepcopy"
	"go.xrstf.de/rudi/pkg/equality"
)

type Semver struct {
	Version blangsemver.Version
}

var (
	_ deepcopy.Copier                  = Semver{}
	_ coalescing.CustomStringCoalescer = Semver{}
	_ equality.Comparer                = Semver{}
)

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
