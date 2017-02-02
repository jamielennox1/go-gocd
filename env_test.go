package gocd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvironment_AddPipeline(t *testing.T) {
	env := NewEnvironment()
	assert.NoError(t, env.AddPipeline("pipeline"))
	assert.Equal(t, len(env.Pipelines), 1)
	for _, p := range env.Pipelines {
		assert.Equal(t, p.Name, "pipeline")
		return
	}
}
func TestEnvironment_AddPipelineExist(t *testing.T) {
	env := NewEnvironment()
	env.AddPipeline("pipeline")
	assert.Error(t, env.AddPipeline("pipeline"))
	assert.Equal(t, len(env.Pipelines), 1)
	for _, p := range env.Pipelines {
		assert.Equal(t, p.Name, "pipeline")
		return
	}
}

func TestEnvironment_DeletePipeline(t *testing.T) {
	env := NewEnvironment()
	env.AddPipeline("pipeline")
	assert.NoError(t, env.DeletePipeline("pipeline"))
	assert.Equal(t, len(env.Pipelines), 0)
}

func TestEnvironment_DeletePipelineNotExist(t *testing.T) {
	env := NewEnvironment()
	env.AddPipeline("pipeline")
	assert.Error(t, env.DeletePipeline("pipeline1"))
	assert.Equal(t, len(env.Pipelines), 1)
}

func TestEnvironment_ExistPipeline(t *testing.T) {
	env := NewEnvironment()
	env.AddPipeline("pipeline")
	assert.Equal(t, env.ExistPipeline("pipeline"), true)
	assert.Equal(t, env.ExistPipeline("pipeline1"), false)
}

func TestEnvironment_AddEnvironmentVariables(t *testing.T) {
	env := NewEnvironment()
	assert.NoError(t, env.AddEnvironmentVariables(&EnvironmentVariable{
		Name:   "test_env",
		Secure: true,
		Crypt:  false,
		Value:  "test"}))
	assert.Equal(t, len(env.EnvironmentVariables), 1)
}

func TestEnvironment_AddEnvironmentVariablesExist(t *testing.T) {
	env := NewEnvironment()
	env.AddEnvironmentVariables(&EnvironmentVariable{
		Name:   "test_env",
		Secure: true,
		Crypt:  false,
		Value:  "test"})
	assert.Error(t, env.AddEnvironmentVariables(&EnvironmentVariable{
		Name:   "test_env",
		Secure: false,
		Crypt:  true,
		Value:  "test1"}))
	assert.Equal(t, len(env.EnvironmentVariables), 1)
}

func TestEnvironment_GetEnvironmentVariablesExist(t *testing.T) {
	env := NewEnvironment()

	envVar1 := EnvironmentVariable{
		Name:   "test_env",
		Secure: true,
		Crypt:  false,
		Value:  "test"}
	env.AddEnvironmentVariables(&envVar1)

	envVar2, err := env.GetEnvironmentVariables(envVar1.Name)

	assert.Equal(t, *envVar2, envVar1)
	assert.Equal(t, err, nil)
}

func TestEnvironment_GetEnvironmentVariablesNotExist(t *testing.T) {
	env := NewEnvironment()

	envVar1 := EnvironmentVariable{
		Name:   "test_env",
		Secure: true,
		Crypt:  false,
		Value:  "test"}
	env.AddEnvironmentVariables(&envVar1)

	envVar2, err := env.GetEnvironmentVariables("test")

	assert.Nil(t, envVar2)
	assert.Error(t, err)
}

func TestEnvironment_DeleteEnvironmentVariables(t *testing.T) {
	env := NewEnvironment()
	env.AddEnvironmentVariables(&EnvironmentVariable{
		Name:   "test_env",
		Secure: true,
		Crypt:  false,
		Value:  "test"})
	assert.NoError(t, env.DeleteEnvironmentVariables("test_env"))
	assert.Equal(t, len(env.EnvironmentVariables), 0)
}

func TestEnvironment_DeleteEnvironmentVariablesNoExist(t *testing.T) {
	env := NewEnvironment()
	env.AddEnvironmentVariables(&EnvironmentVariable{
		Name:   "test_env",
		Secure: true,
		Crypt:  false,
		Value:  "test"})
	assert.Error(t, env.DeleteEnvironmentVariables("test"))
	assert.Equal(t, len(env.EnvironmentVariables), 1)
}
