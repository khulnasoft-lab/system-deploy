package actions

import (
	"fmt"
	"testing"

	"github.com/khulnasoft-lab/system-conf/conf"
	"github.com/khulnasoft-lab/system-deploy/pkg/deploy"
	"github.com/stretchr/testify/assert"
)

func TestRegisterPlugin(t *testing.T) {
	setupCalled := false

	plg := Plugin{
		Name: "Test",
		Setup: func(t deploy.Task, opts conf.Section) (Action, error) {
			if setupCalled {
				return nil, fmt.Errorf("called")
			}

			setupCalled = true
			return nil, nil
		},
	}

	assert.NoError(t, Register(plg))
	assert.Error(t, Register(plg))

	p, ok := GetPlugin("test")
	assert.True(t, ok)
	assert.Equal(t, plg.Name, p.Name)

	_, ok = GetPlugin("unknown")
	assert.False(t, ok)

	actionList := ListActions()
	assert.Len(t, actionList, 1)

	assert.False(t, setupCalled)
	_, err := Setup("test", nil, deploy.Task{}, conf.Section{})
	assert.Error(t, err)
	assert.True(t, setupCalled)

	_, err = Setup("test", nil, deploy.Task{}, conf.Section{})
	assert.Error(t, err)
	assert.True(t, setupCalled)

	actions = make(map[string]*Plugin)
}
