# Sequence Diagrams

Visual representation of key message flows in mautrix-viber.

## Viber → Matrix Message Flow

```
Viber Platform          Bridge              Database            Matrix
     |                    |                    |                  |
     |--webhook---------->|                    |                  |
     |  POST /webhook     |                    |                  |
     |                    |                    |                  |
     |                    |--verify sig------->|                  |
     |                    |                    |                  |
     |                    |--parse payload     |                  |
     |                    |                    |                  |
     |                    |--store sender----->|                  |
     |                    |                    |                  |
     |                    |--get room--------->|                  |
     |                    |<--room ID----------|                  |
     |                    |                    |                  |
     |                    |--------forward msg------------------>|
     |                    |                    |                  |
     |                    |--store mapping---->|                  |
     |                    |                    |                  |
     |<--200 OK-----------|                    |                  |
```

## Matrix → Viber Message Flow

```
Matrix              Bridge              Database            Viber API
  |                    |                    |                  |
  |--event------------>|                    |                  |
  |  MessageEvent      |                    |                  |
  |                    |                    |                  |
  |                    |--get mapping------>|                  |
  |                    |<--viber chat ID----|                  |
  |                    |                    |                  |
  |                    |--format message    |                  |
  |                    |                    |                  |
  |                    |--------send msg---------------------->|
  |                    |                    |                  |
  |                    |--store mapping---->|                  |
  |                    |                    |                  |
  |                    |<--response---------|                  |
```

## Webhook Registration Flow

```
Bridge                Viber API
  |                      |
  |--POST /set_webhook-->|
  |  {url, events}        |
  |                      |
  |<--200 OK-------------|
  |  {status: 0}         |
  |                      |
  |--verify response     |
  |                      |
```

## Health Check Flow

```
Orchestrator          Bridge              Database
  |                      |                    |
  |--GET /healthz------->|                    |
  |                      |                    |
  |<--200 OK-------------|                    |
  |                      |                    |
  |--GET /readyz-------->|                    |
  |                      |                    |
  |                      |--ping------------->|
  |                      |<--OK---------------|
  |                      |                    |
  |<--200 OK-------------|                    |
```

## Error Handling Flow

```
Client                Bridge              External API
  |                      |                      |
  |--request------------>|                      |
  |                      |                      |
  |                      |--API call---------->|
  |                      |                      |
  |                      |<--error--------------|
  |                      |                      |
  |                      |--retry (backoff)    |
  |                      |                      |
  |                      |--API call---------->|
  |                      |                      |
  |                      |<--success-----------|
  |                      |                      |
  |<--200 OK-------------|                      |
```

