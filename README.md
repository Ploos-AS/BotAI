# BotAI

BotAI is the optional shared AI service for Ploos IRC bots.

## Status

**M0.7 — runtime expert reload and status implemented; CI qualification pending.**

M0 provides a small Go service with a provider abstraction, initial expert registry
(`general`, `irc`, `amiga`), bounded JSON API, deterministic test provider,
unit tests, GitHub Actions CI and an Alpine OCI image.

## Run

```sh
go run ./cmd/botai
```

BotAI listens on `127.0.0.1:8090` by default. Override with `BOTAI_LISTEN`.

Endpoints:

- `GET /healthz`
- `GET /v1/experts`
- `POST /v1/chat`

Example:

```sh
curl -s http://127.0.0.1:8090/v1/chat \
  -H 'Content-Type: application/json' \
  -d '{"expert":"irc","message":"What is IRC 433?"}'
```

The M0 echo provider is deliberately deterministic and has no network dependency.
M0.1 adds an optional OpenAI-compatible provider. Select it with:

```sh
BOTAI_PROVIDER=openai-compatible \
BOTAI_BASE_URL=http://127.0.0.1:11434/v1 \
BOTAI_MODEL=my-model \
go run ./cmd/botai
```

`BOTAI_API_KEY` is optional for endpoints that do not require authentication.
The deterministic `echo` provider remains the default and is used for offline qualification.

## Standalone-first rule

BotAI is an optional enhancement, never part of the IRC bot core. LuCa, Engo, AmBot, and other participating bots MUST remain fully functional IRC bots without BotAI, BotWeb, PBMP, or any other external Ploos service.

BotAI integrations MUST fail open with respect to ordinary bot operation: if BotAI is disabled, unavailable, misconfigured, or removed, only AI-specific capabilities may become unavailable. IRC connectivity, lifecycle, built-in commands, scripting/plugin systems, ACLs, persistence, and other non-AI bot functionality MUST continue independently.

BotAI MAY provide conversational AI, expert roles, model routing, context and knowledge/RAG services through optional adapters. It MUST NOT require BotWeb. Management through PBMP/BotWeb may be supported, but direct bot-to-BotAI integration must remain possible where appropriate.

## Integration principle

```text
IRC bot core (standalone)
    |
    +-- optional PBMP adapter --> BotWeb
    |
    +-- optional AI adapter ---> BotAI
```

This keeps AI and web management independently optional and prevents hidden service dependencies.

See `docs/M0.md` through `docs/M0.7.md` and `docs/ARCHITECTURE.md`.
