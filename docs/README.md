# Relay Documentation

Welcome to the Relay documentation. This directory contains detailed guides and reference material for the Universal LLM Provider Runtime.

## Documentation Index

### Core Documentation

- **[Architecture](ARCHITECTURE.md)** — System architecture, data flow, and design decisions
- **[Host-Engine Boundary](architecture/HOST-ENGINE-BOUNDARY.md)** — Stable `engine/llm/graph/tools` contract
- **[Relay Enterprise](design/RELAY-ENTERPRISE.md)** — Enterprise surfaces
- **[Provider Setup Guide](guides/CREDENTIAL-SETUP-FLOW.md)** — How to configure credentials and providers
- **[Dynamic Model Discovery](guides/DYNAMIC-MODEL-DISCOVERY.md)** — Architecture and implementation details for live model discovery
- **[Decentralized Relay routing](architecture/DECENTRALIZED-RELAY.md)** — Instance-local routing and signed peer manifests
- **[Feature-oriented architecture](architecture/FEATURE-MONOREPO.md)** — Package layout and layering rules
- **[OpenAPI](../api/openapi.yaml)** — Contract for the internal HTTP server in `internal/api` (health, prompt, nodes, aliases, analytics); `/v1/chat/completions`, `/rerank` and `/ready` are served but not yet in the spec. relay ships no binary that starts this server.

### Quick Links

- **[README](../README.md)** — Project overview and quick start
- **[Contributing Guide](../CONTRIBUTING.md)** — How to contribute to Relay
- **[Security Policy](../SECURITY.md)** — Security reporting and best practices
- **[Changelog](../CHANGELOG.md)** — Version history and release notes

### Examples

The [`examples/`](../examples/) directory contains runnable code samples:

- **Basic Chat** — Simple synchronous chat with a single provider
- **Streaming** — Server-sent events streaming with continuation
- **Multi-Provider** — Using fallback chains across multiple providers

## Documentation Structure

```
docs/
├── README.md                          # This file
├── ARCHITECTURE.md                    # System architecture
├── architecture/
│   ├── HOST-ENGINE-BOUNDARY.md        # Host contract, frozen engine-internal types
│   ├── DECENTRALIZED-RELAY.md          # Instance-local routing, signed peer manifests
│   └── FEATURE-MONOREPO.md            # Package layout and layering
├── design/
│   └── RELAY-ENTERPRISE.md             # Enterprise surfaces (design)
├── guides/
│   ├── CREDENTIAL-SETUP-FLOW.md
│   └── DYNAMIC-MODEL-DISCOVERY.md
└── plans/                             # Remediation plans and reviews
```

The HTTP contract lives outside this directory, at
[`../api/openapi.yaml`](../api/openapi.yaml). Retry/fallback, caching and
routing strategies are described in [ARCHITECTURE.md](ARCHITECTURE.md); there
are no separate guides for them yet.

## For Developers

If you're contributing to Relay:

1. Read [CONTRIBUTING.md](../CONTRIBUTING.md) for development setup
2. Review [ARCHITECTURE.md](ARCHITECTURE.md) to understand the system
3. Check [AGENTS.md](../AGENTS.md) for AI agent context and conventions
4. Run `make ci` locally before submitting PRs

## For Users

If you're using Relay in your application:

1. Start with the [Quick Start](../README.md#quick-start) in the main README
2. Review the [Usage examples](../README.md#usage) for common patterns
3. Check the [examples/](../examples/) directory for complete code samples
4. Read the [Provider Setup Guide](guides/CREDENTIAL-SETUP-FLOW.md) for credential configuration

## API Reference

API documentation is available at:
- **[pkg.go.dev](https://pkg.go.dev/github.com/GrayCodeAI/relay)** — Generated Go documentation
- **[ARCHITECTURE.md](ARCHITECTURE.md)** — Core abstractions and interfaces

## Support

- **Issues**: [GitHub Issues](https://github.com/GrayCodeAI/relay/issues)
- **Security**: See [SECURITY.md](../SECURITY.md) for vulnerability reporting
