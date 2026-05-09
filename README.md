# Factory CLI

Factory is a Go + Cobra CLI for creating and running configurable agent command workflows.

This repository root is the Factory app itself. The app is the CLI that lets users create, version, observe, test, improve, and run reusable workflows and steps.

The `.factory/` directory in this repository is not the app. It is this repository's Factory workspace: a place to store command mappings, workflow definitions, handoffs, artifacts, events, logs, and other runtime/project state while we dogfood Factory on itself.

## Current Model

Factory command behavior is materialized as prompt-program files and mapped through a structured registry in a workspace:

```text
.factory/
├── 00-control-room/
│   ├── a-config/
│   │   ├── commands.json
│   │   ├── event-bindings.json
│   │   └── factory.json
│   └── b-commands/
│       └── <command>.md
```

The source of truth is:
- command prompt files for executable pseudocode
- `commands.json` for command discovery and metadata
- `event-bindings.json` for autonomous continuation rules
- event logs for runtime history

The CLI must not hard-code workflow commands. Users create commands and bindings to define their own production line.

Factory supports two execution styles:
- Manual: the user runs one command or next step at a time.
- Autonomous: the CLI follows event bindings and continues through the configured production line until completion, blockage, approval requirement, or an iteration guard.

## Basic Usage

Initialize a workspace:

```bash
go run ./cmd/factory init
```

Create and register a command:

```bash
go run ./cmd/factory command create load-intake \
  --alias load \
  --description "Load one intake item into the factory."
```

Create an event binding:

```bash
go run ./cmd/factory event bind create \
  --where command=load-intake \
  --where result.status=succeeded \
  --run drain-work-package \
  --input '{"package_path":"$result.data.package_path"}'
```

Run a command:

```bash
go run ./cmd/factory run load-intake --data '{"text":"..."}'
```

Visualize the current factory:

```bash
go run ./cmd/factory visualize
go run ./cmd/factory visualize --format json
go run ./cmd/factory visualize --output .factory/00-control-room/e-state/factory-flow.mmd
```

The visualization reads the configured workspace lanes, registered commands, event bindings, sectors, sector actions, and prompt references to sectors.

## Development

Run tests with a writable Go cache:

```bash
env GOCACHE=/tmp/go-build GOMODCACHE=/tmp/go-mod go test ./...
```

Build the CLI:

```bash
env GOCACHE=/tmp/go-build GOMODCACHE=/tmp/go-mod go build -o /tmp/factory-cli ./cmd/factory
```

## Product Triangulation

The current product baseline is the gamified factory model:

- [Triangulation](docs/triangulation.md)
- [Roadmap](docs/roadmap.md)
- [Use Cases](docs/use-cases.md)
- [Architecture](docs/architecture.md)
