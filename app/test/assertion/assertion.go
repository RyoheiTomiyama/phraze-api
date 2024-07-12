package assertion

import (
	"testing"

	"github.com/RyoheiTomiyama/phraze-api/util/errutil"
	"github.com/stretchr/testify/assert"
)

func AssertError(t *testing.T, expectedMsg string, expectedCode errutil.ErrorCode, actualErr error) {
	t.Helper()

	assert.Error(t, actualErr)

	customErr, ok := actualErr.(errutil.IError)
	if !ok {
		t.Fatalf("err is not *errutil.Error: %#v", actualErr)
	}

	assert.Equal(t, expectedMsg, customErr.Error())
	assert.Equal(t, int(expectedCode), customErr.Code())
}
