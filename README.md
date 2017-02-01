# endpoint
A robust HTTP request router for Go.

Example:
```
e := endpoint.New()
e.Match("id", `^\d+$`)
e.Route("GET /", home)
e.Route("POST /note", noteCreate)
e.Route("GET /note/:id", noteRetrieve)
e.Serve("127.0.0.1:3333")
```
