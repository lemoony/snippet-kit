package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"

	mocks "github.com/lemoony/snipkit/mocks/app"
)

func Test_ManagerSync(t *testing.T) {
	app := mocks.App{}
	app.On("SyncManager").Return(nil, nil)

	err := runMockedTest(t, []string{"manager", "sync"}, withApp(&app))

	assert.NoError(t, err)
	app.AssertNumberOfCalls(t, "SyncManager", 1)
}
