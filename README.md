# BotAI

BotAI is the optional shared AI service for Ploos IRC bots.

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
