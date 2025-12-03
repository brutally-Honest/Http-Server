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