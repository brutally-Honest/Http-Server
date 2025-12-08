![alt text](https://media.npr.org/assets/img/2023/05/26/honest-work-meme-c7034f8bd7b11467e1bfbe14b87a5f6a14a5274b.jpg?s=800&c=85&f=webp)

1. Unless it's HTTP/1.0, an HTTP/1.1 server closes the connection (mostly) because:
    * Client asks for it:
        - "Connection: close"
    * Server enforces a connection policy:
        - Max requests per connection reached
    * Server decides to close:
        - Idle-timeout
        - Graceful shutdown / reload
        - Memory or FD pressure
2. HTTP 1.1 supports pipelining, sequencing also, but the default behaviour of most modern clients is to close the connection after a single req–res cycle
3. To test if requests reuse a connection using curl, use `--next` for multiple URLs.  
This triggers HTTP/1.1 keep-alive sequencing (sequential request–response cycles on the same TCP connection).
4. Transfer-encoding chunked in HTTP/1.1 is ~28 years old (HTTP/1.1 -> 1997, today 2025), almost three decades, still alive and kicking
5. Lapse in understanding: Streaming ≠ SSE!
Its a protocol on top of HTTP which deals with UTF-8 only and requires a very specific event-stream format (headers + body)
6. Request can only be cancelled by server with connection still alive, its only connection close from the client
