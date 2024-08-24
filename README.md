# Gin Middleware Extension

Gin-Middleware is a middleware extension of the [Gin](https://github.com/gin-gonic/gin) framework in [Go](https://go.dev/), including common functions such as guards, pipes, interceptors, etc.

**key features are:**

- guard
  - [AuthGuard](./guard.go#L23)
- pipe
  - [QueryPipe](./pipe.go#L14)
  - [BodyPipe](./pipe.go#L31)
- alarm
  - [FeishuWebHookAlarm](./alarm.go#L16)

## Getting started

### Prerequisites

Gin requires [Go](https://go.dev/) version [1.21](https://go.dev/doc/devel/release#go1.21.0) or above.

### Getting Gin-Middleware

With [Go's module support](https://go.dev/wiki/Modules#how-to-use-modules), `go [build|run|test]` automatically fetches the necessary dependencies when you add the import in your code:

```sh
import "github.com/zhangyuan722/gin-middleware"
```

Alternatively, use `go get`:

```sh
go get -u github.com/zhangyuan722/gin-middleware@latest
```

### Test
```sh
go test -v
```

