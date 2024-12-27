# Tome Dome
Powers https://dota.tomedome.io

Pave a path to Immortal by eliminating game knowledge gaps!


# Setup
* Install `docker` and `docker-compose` for building and running
* Install `npm` for running js unit tests
* `make build` to build container image and install js test toolchain


# Run tests
* `make test` for api and js unit tests
* `make integration-test` runs api integration test against public Stratz API (requires tomedome_STRATZ_API_KEY env var set)


# Run app
1. `make run-server`
2. visit http://localhost:8080


# Architecture

```mermaid
graph TD
  subgraph GCP
    direction TB
    lb["L7 LB"]
    bucket["static assets (GCS bucket)"]
    api["rest api (cloud func)"]

    %% connections

    lb --> bucket
    lb --> api
    dns1 --> lb
    dns2 --> lb

    subgraph dns["cloud dns"]
      dns1["dota.tomedome.io"]
      dns2["api.tomedome.io"]
    end
  end
  
  user[user browser] -->|fetch static content| dns1
  user -->|fetches json data| dns2
 ```


# todo
- [ ] ci/cd pipeline
- [ ] google analytics


# Roadmap

Current state is MVP for "shipping the smallest useful thing"

* Quiz behavior:
    * Have fixed number of questions
    * Give you a score at the end
    * Multiple choice
    * Keep track of quiz results
* Content:
    * Ability stats (e.g., "how much movement speed does Scorched Earth give at level 3?")    
    * Items
    * Synergies ("You're in a lane as PA with CM against Tide and Venge. What items should you pick?")