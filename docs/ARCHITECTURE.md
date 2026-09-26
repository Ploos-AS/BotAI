# BotAI architecture

## M0

BotAI is an out-of-process, optional AI capability service for IRC bots.

The bot core is always authoritative for IRC connectivity, lifecycle, ACLs,
commands, scripting and persistence. BotAI receives bounded requests through an
adapter and returns AI-specific responses. A BotAI outage may disable AI
features, but MUST NOT stop or degrade ordinary IRC bot operation.

M0 layers:

- HTTP API: health, expert discovery and chat request
- expert registry: named profiles independent of model providers
- provider interface: model/backend abstraction
- deterministic echo provider: M0 qualification without external dependencies

Future providers may include local and cloud model APIs. PBMP/BotWeb management
is optional and is not in the BotAI request path.
