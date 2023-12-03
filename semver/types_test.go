// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package semver

import (
	"testing"

	blangsemver "github.com/blang/semver/v4"
)

func TestSemverDeepCopy(t *testing.T) {
	input := blangsemver.MustParse("1.0.0-beta.1")
	sv := Semver{Version: input}

	copied, err := sv.DeepCopy()
	if err != nil {
		t.Fatalf("Failed to deepcopy Semver object: %v", err)
	}

	copiedSV, ok := copied.(Semver)
	if !ok {
		t.Fatalf("DeepCopy did not return Semver, but %T", copied)
	}

	if !input.Equals(copiedSV.Version) {
		t.Fatalf("Expected %s, but copy is %s", input.String(), copiedSV.Version.String())
	}

	// change the copy
	copiedSV.Version.Pre[0].VersionStr = "foo"

	if input.Equals(copiedSV.Version) {
		t.Fatal("Expected to only change copy, but changed original, too.")
	}
}
