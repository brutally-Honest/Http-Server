# HTTP Server from Scratch (Go)

Raw TCP → HTTP server (HTTP/1.0 and HTTP/1.1)  
No `net/http`

This project is **not a framework or a showcase**.  
It exists to remove ambiguity about how HTTP actually works under the hood.

---

## Why this exists

Most HTTP knowledge stays abstract because frameworks hide everything:
buffers, partial reads, connection lifecycles, streaming, timeouts.

This server was built to understand:
- how bytes move over the wire
- how HTTP is framed on top of TCP
- how servers actually parse, stream, and respond
- where protocol boundaries really are

Correctness and clarity mattered more than completeness.

---

## What I learned

### TCP & networking fundamentals
- TCP sockets:
  - listening socket vs accepted connection socket
- TCP is a **bidirectional byte stream**, not message-based
- Partial reads and writes are normal
- TCP guarantees:
  - ordering
  - retransmission
  - reliability
- Connection lifecycle:
  - 3-way handshake (SYN → SYN-ACK → ACK)
  - 4-way teardown (FIN / ACK sequence)

---

### HTTP protocol internals
- Manual parsing of HTTP requests:
  - request line
  - headers
  - CRLF framing
  - body
- HTTP/1.0 vs HTTP/1.1 differences
- Persistent connections and reuse
- HTTP pipelining (spec support vs browser reality)
- Header handling:
  - `Content-Length`
  - `Transfer-Encoding: chunked`
  - `Content-Type`
- Chunked transfer encoding:
  - hex size
  - CRLF framing
- Streaming semantics:
  - client → server
  - server → client
- SSE (Server-Sent Events):
  - built on HTTP streaming
  - UTF-8 only
  - text/event-stream framing

---

### Server-side behavior
- Incremental request parsing from a TCP stream
- Streaming responses using chunked encoding
- Error handling for malformed input
- Connection reuse across multiple request–response cycles

---

### Limits & timeouts (config-driven)
- Request size limits
- Header size limits
- Body size limits
- Read timeouts
- Write timeouts

These are enforced explicitly and not hardcoded.

---

### HTTP/1.1 compliance handling
- Mandatory `Host` header enforcement
- Rejection of multiple `Content-Length` headers
- Clear separation of:
  - client-initiated close
  - server-initiated close (timeouts, limits, shutdown)

---

### Routing & data structures
- Radix tree–based router
- Static, parameter, and wildcard routes
- Applied DSA concepts in a real system (not just problems)

---

### Go-specific learnings
- Context usage:
  - cancellation
  - request lifecycle propagation
- Error handling across streaming boundaries
- Writing network code without `net/http`

---

### Tooling
- Used `curl` to:
  - inspect raw HTTP traffic
  - test streaming behavior
  - debug headers and chunked responses

---

## Project structure
```
├── cmd/
│   └── server/
│       └── main.go
└── internal/
    ├── config/
    │   └── config.go
    ├── server/
    │   ├── server.go
    │   └── connection.go
    ├── router/
    │   └── router.go
    ├── request/
    │   ├── request.go
    │   ├── headers.go
    │   ├── body.go
    │   ├── readers.go
    │   ├── parser.go
    │   ├── validators.go
    │   ├── line.go
    │   └── errors.go
    └── response/
        ├── response.go
        ├── headers.go
        ├── writers.go
        ├── flush.go
        ├── chunked.go
        └── errors.go
```
---

## Design choices & scope

### Implemented intentionally
- Configurable limits and timeouts
- Manual HTTP parsing and validation
- Streaming support in both directions
- Radix-tree routing
- Connection reuse with explicit constraints

### Out of scope (by design)
- TLS
- HTTP/2 or HTTP/3
- Compression (gzip, brotli)
- Trailer headers
- `Expect: 100-continue`
- Middleware chaining
- JSON handling
- File uploads and downloads

This project favors understanding protocol mechanics over RFC completeness.

---

## Potential improvements

- Graceful shutdown with signal fan-out
- Explicit connection and request state machines
- Parser fuzz testing
- Deeper Go concurrency patterns (channels, worker models)

---

## Non-goals

- Production readiness
- Feature parity with `net/http`
- Performance optimization
- Backward compatibility with broken clients

---

## Final note

This project was built to **see the edges**, not to hide them.

If something feels “manual” or “verbose”, that is intentional —  
because that’s where the real learning happens.
