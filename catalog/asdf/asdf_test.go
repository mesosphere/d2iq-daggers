// Copyright 2022 D2iQ, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package asdf

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTraverseNonCommentLines(t *testing.T) {
	content := `
foo
# foo
## foo
	 # foo
`

	buf := bytes.NewBufferString(content)

	count := 0
	traverseNonCommentLines(buf, func(line string) {
		if strings.Contains(line, "foo") {
			count += 1
		}
	})

	assert.Equal(t, count, 1)
}

func TestPluginVersions_GetVersionOrDefault(t *testing.T) {
	versions := PluginVersions{
		"foo": {
			Version:       "1.2.3",
			VersionFreeze: true,
		},
	}

	assert.Equal(t, versions.GetVersionOrDefault("foo", "v", "latest"), "v1.2.3")
	assert.Equal(t, versions.GetVersionOrDefault("foo", "", "latest"), "1.2.3")
	assert.Equal(t, versions.GetVersionOrDefault("bar", "v", "1.0.0"), "1.0.0")
}
