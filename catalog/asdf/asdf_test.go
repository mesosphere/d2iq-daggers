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
