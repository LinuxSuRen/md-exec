package cli

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScriptRunners(t *testing.T) {
	runners := ScriptRunners{}
	assert.Equal(t, 0, runners.Size())
	assert.Nil(t, runners.GetRunner("fake"))
	assert.Equal(t, []string{}, runners.GetTitles())

	runners = NewScriptRunners()
	quitRunner := runners.GetRunner("Quit")
	if assert.NotNil(t, quitRunner) {
		assert.Equal(t, "Quit", quitRunner.GetTitle())
		assert.Nil(t, quitRunner.Run(context.Background()))
	}
	assert.Equal(t, []string{"Quit"}, runners.GetTitles())
}
