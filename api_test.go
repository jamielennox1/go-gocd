package gocd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createPath(path string) string {
	return fmt.Sprintf("./test_data/%s.json", path)
}

func getEnvironmentHandler(w http.ResponseWriter, r *http.Request) {
	if strings.Compare(r.Method, "GET") != 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
		fmt.Fprint(w, fmt.Sprintf(`{"Error":"method %s != GET"}`, r.Method))
		return
	}
	data, err := ioutil.ReadFile(createPath("get_environment"))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
		fmt.Fprint(w, fmt.Sprintf(`{"Error":"%v"}`, err))
		return
	}
	w.Header().Set("ETag", "123456789")
	w.Write(data)
}

func getEnvironmentsHandler(w http.ResponseWriter, r *http.Request) {
	if strings.Compare(r.Method, "GET") != 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
		fmt.Fprint(w, fmt.Sprintf(`{"Error":"method %s != GET"}`, r.Method))
		return
	}
	data, err := ioutil.ReadFile(createPath("get_environments"))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
		fmt.Fprint(w, fmt.Sprintf(`{"Error":"%v"}`, err))
		return
	}
	w.Header().Set("ETag", "123456789")
	w.Write(data)
}

func TestClient_GetEnvironment(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(getEnvironmentHandler))
	defer server.Close()

	client := New(server.URL, "", "")
	if env, err := client.GetEnvironment("TEST"); err != nil {
		t.Error(err)
		t.Fail()
	} else {
		assert.Equal(t, env.Name, "TEST")
	}
}

func TestClient_GetEnvironments(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(getEnvironmentsHandler))
	defer server.Close()

	client := New(server.URL, "", "")
	if envs, err := client.GetEnvironments(); err != nil {
		t.Error(err)
		t.Fail()
	} else {
		assert.Equal(t, len(envs.Embeded.Environments), 1)

		for _, env := range envs.Embeded.Environments {
			assert.Equal(t, len(env.Agents), 1)
			assert.Equal(t, len(env.EnvironmentVariables), 2)
			assert.Equal(t, len(env.Pipelines), 1)
			break
		}
	}
}
