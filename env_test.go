package gocd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvironment_AddPipeline(t *testing.T) {
	env := Environment{}
	assert.Equal(t, env.AddPipeline("pipeline"), true)
	assert.Equal(t, len(env.Pipelines), 1)
	for _, p := range env.Pipelines {
		assert.Equal(t, p.Name, "pipeline")
	}
}
func TestEnvironment_AddPipelineExist(t *testing.T) {
	env := Environment{Pipelines: []ShortPipeline{ShortPipeline{Name: "pipeline"}}}
	assert.Equal(t, env.AddPipeline("pipeline"), false)
	assert.Equal(t, len(env.Pipelines), 1)
	for _, p := range env.Pipelines {
		assert.Equal(t, p.Name, "pipeline")
	}
}

func TestEnvironment_DeletePipeline(t *testing.T) {
	env := Environment{Pipelines: []ShortPipeline{ShortPipeline{Name: "pipeline"}}}
	assert.Equal(t, env.DeletePipeline("pipeline"), true)
	assert.Equal(t, len(env.Pipelines), 0)
}

func TestEnvironment_DeletePipelineNotExist(t *testing.T) {
	env := Environment{Pipelines: []ShortPipeline{ShortPipeline{Name: "pipeline"}}}
	assert.Equal(t, env.DeletePipeline("pipeline1"), false)
	assert.Equal(t, len(env.Pipelines), 1)
}

func TestEnvironment_ExistPipeline(t *testing.T) {
	env := Environment{Pipelines: []ShortPipeline{ShortPipeline{Name: "pipeline"}}}
	assert.Equal(t, env.ExistPipeline("pipeline"), true)
	assert.Equal(t, env.ExistPipeline("pipeline1"), false)
}
