# GoCD Client

Go Lang library to access [GoCD API](https://api.gocd.io/16.7.0/)

## Usage
```go
package main

import (
  "github.com/mhanygin/go-gocd"
)

func main() {
  client := gocd.New("http://gocd.com:8153", "login", "password")
  // ... do whatever you want with the client
}

## API Endpoints Pending
- Agents
  - [ ] Get all Agents
  - [ ] Get one Agent
  - [ ] Update an Agent
  - [ ] Disable Agent
  - [ ] Delete an Agent
  - [ ] Agent job run history
- Users
  - [ ] Get all Users
  - [ ] Get one user
  - [ ] Create a user
  - [ ] Update a user
  - [ ] Delete a user
- Materials
  - [ ] Get all Materials
  - [ ] Get material modifications
  - [ ] Notify SVN materials
  - [ ] Notify git materials
- Backups
  - [ ] Create a backup
- Pipeline Group
  - [x] Config listing
- Artifacts
  - [ ] Get all Artifacts
  - [ ] Get artifact file
  - [ ] Get artifact directory
  - [ ] Create artifact
  - [ ] Append to artifact
- Pipelines
  - [x] Get pipeline instance
  - [ ] Get pipeline status
  - [x] Pause a pipeline
  - [x] Unpause a pipeline
  - [ ] Releasing a pipeline lock
  - [x] Scheduling Pipelines
- Stages
  - [ ] Cancel Stage
  - [ ] Get Stage instance
  - [ ] Get stage history
- Jobs
  - [ ] Get Scheduled Jobs
  - [ ] Get Job history
- Properties
  - [ ] Get all job Properties
  - [ ] Get one property
  - [ ] Get historical properties
  - [ ] Create property
- Configurations
  - [ ] List all modifications
  - [ ] Get repository modification diff
  - [ ] Get Configuration  
- Environment Config
  - [x] Get all environments
  - [x] Get environment config
  - [x] Create an environment
  - [x] Update an environment
  - [x] Delete an environment
- [ ] Dashboard
  - [ ] Get Dashboard
- Pipeline Config
  - [x] Get pipeline Configuration
  - [x] Edit Pipeline configuration
  - [x] Create Pipeline
  - [x] Delete Pipeline