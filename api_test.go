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

func TestClient_GetEnvironment(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Compare(r.Method, "GET") != 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
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
		w.Header().Set("Etag", "123456789")
		w.Write(data)
	}))
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
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Compare(r.Method, "GET") != 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
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
		w.Header().Set("Etag", "123456789")
		w.Write(data)
	}))
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

func TestClient_PausePipeline(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Compare(r.Method, "POST") != 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprint(w, fmt.Sprintf(`{"Error":"method %s != POST"}`, r.Method))
			return
		}
		if strings.Compare(r.Header.Get("Confirm"), "true") != 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, fmt.Sprint(`{"Error":"header Confirm != true"}`))
			return
		}
		defer r.Body.Close()
		if body, err := ioutil.ReadAll(r.Body); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, fmt.Sprintf(`{"Error":"%v"}`, err))
			return
		} else if strings.Compare(string(body), "pauseCause=take some rest") != 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, fmt.Sprintf(`{"Error":"Body %s != pauseCause=take some rest"}`, body))
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := New(server.URL, "", "")
	if err := client.PausePipeline("pipeline"); err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestClient_UnpausePipeline(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Compare(r.Method, "POST") != 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprint(w, fmt.Sprintf(`{"Error":"method %s != POST"}`, r.Method))
			return
		}
		if strings.Compare(r.Header.Get("Confirm"), "true") != 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, fmt.Sprint(`{"Error":"header Confirm != true"}`))
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := New(server.URL, "", "")
	if err := client.UnpausePipeline("pipeline"); err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestClient_SchedulePipeline(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Compare(r.Method, "POST") != 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprint(w, fmt.Sprintf(`{"Error":"method %s != POST"}`, r.Method))
			return
		}
		if strings.Compare(r.Header.Get("Confirm"), "true") != 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, fmt.Sprint(`{"Error":"header Confirm != true"}`))
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprint(w, "Request to schedule pipeline pipeline accepted")
	}))
	defer server.Close()

	client := New(server.URL, "", "")
	if err := client.SchedulePipeline("pipeline", []byte{}); err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestClient_GetGroups(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Compare(r.Method, "GET") != 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprint(w, fmt.Sprintf(`{"Error":"method %s != GET"}`, r.Method))
			return
		}

		data, err := ioutil.ReadFile(createPath("get_groups"))
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNoContent)
			fmt.Fprint(w, fmt.Sprintf(`{"Error":"%v"}`, err))
			return
		}
		w.Header().Set("Etag", "123456789")
		w.Write(data)
	}))
	defer server.Close()

	client := New(server.URL, "", "")
	if grps, err := client.GetGroups(); err != nil {
		t.Error(err)
		t.Fail()
	} else {
		assert.Equal(t, len(*grps), 2)
	}
}

func TestClient_GetPipelineConfig(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Compare(r.Method, "GET") != 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprint(w, fmt.Sprintf(`{"Error":"method %s != GET"}`, r.Method))
			return
		}

		data, err := ioutil.ReadFile(createPath("get_pipeline_config"))
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNoContent)
			fmt.Fprint(w, fmt.Sprintf(`{"Error":"%v"}`, err))
			return
		}
		w.Header().Set("Etag", "123456789")
		w.Write(data)
	}))
	defer server.Close()

	client := New(server.URL, "", "")
	if pipeline, err := client.GetPipelineConfig("my_pipeline"); err != nil {
		t.Error(err)
		t.Fail()
	} else {
		assert.Equal(t, pipeline.Name, "my_pipeline")
	}
}

//func TestClient_SetPipelineConfig(t *testing.T) {
//	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		if strings.Compare(r.Method, "PUT") != 0 {
//			w.Header().Set("Content-Type", "application/json")
//			w.WriteHeader(http.StatusMethodNotAllowed)
//			fmt.Fprint(w, fmt.Sprintf(`{"Error":"method %s != PUT"}`, r.Method))
//			return
//		}
//		data, err := ioutil.ReadFile(createPath("post_pipeline_config"))
//		if err != nil {
//			w.Header().Set("Content-Type", "application/json")
//			w.WriteHeader(http.StatusNoContent)
//			fmt.Fprint(w, fmt.Sprintf(`{"Error":"%v"}`, err))
//			return
//		}
//		val1 := make(map[string]interface{})
//		if err := json.Unmarshal(data, &val1); err != nil {
//			w.Header().Set("Content-Type", "application/json")
//			w.WriteHeader(http.StatusNoContent)
//			fmt.Fprint(w, fmt.Sprintf(`{"Error":"%v"}`, err))
//			return
//		}
//
//		defer r.Body.Close()
//		if body, err := ioutil.ReadAll(r.Body); err != nil {
//			w.Header().Set("Content-Type", "application/json")
//			w.WriteHeader(http.StatusBadRequest)
//			fmt.Fprint(w, fmt.Sprintf(`{"Error":"%v"}`, err))
//			return
//		} else {
//			val2 := make(map[string]interface{})
//			if err := json.Unmarshal(body, &val2); err != nil {
//				w.Header().Set("Content-Type", "application/json")
//				w.WriteHeader(http.StatusNoContent)
//				fmt.Fprint(w, fmt.Sprintf(`{"Error":"%v"}`, err))
//				return
//			}
//
//			if !reflect.DeepEqual(val1["stages"], val2["stages"]) {
//				w.Header().Set("Content-Type", "application/json")
//				w.WriteHeader(http.StatusNoContent)
//				fmt.Fprint(w, fmt.Sprintf(`{"Error": "bad content"}`))
//			} else {
//				w.WriteHeader(http.StatusOK)
//			}
//		}
//	}))
//	defer server.Close()
//
//	client := New(server.URL, "", "")
//
//	data, err := ioutil.ReadFile(createPath("post_pipeline_config"))
//	if err != nil {
//		t.Error(err)
//		t.Fail()
//	}
//
//	pipeline := PipelineConfig{}
//
//	if err := json.Unmarshal(data, &pipeline); err != nil {
//		t.Error(err)
//		t.Fail()
//	}
//
//	if err := client.SetPipelineConfig(&pipeline); err != nil {
//		t.Error(err)
//		t.Fail()
//	}
//}
