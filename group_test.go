package gocd

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGroup_Exist(t *testing.T) {
	data := `{
    "pipelines": [
      {
        "stages": [
          {
            "name": "up42_stage"
          }
        ],
        "name": "pp3",
        "materials": [
          {
            "description": "URL: https://github.com/gocd/gocd, Branch: master",
            "fingerprint": "2d05446cd52a998fe3afd840fc2c46b7c7e421051f0209c7f619c95bedc28b88",
            "type": "Git"
          }
        ],
        "label": "${COUNT}"
      }
    ],
    "name": "second"
  }`
	grp := Group{}
	if err := json.Unmarshal([]byte(data), &grp); err != nil {
		t.Error(err)
		t.Fail()
	}
	assert.Equal(t, grp.Exist("pp3"), true)
	assert.Equal(t, grp.Exist("pp"), false)
}
