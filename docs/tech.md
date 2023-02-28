请求的流程

服务器接受一个来自前端的请求以后：

- 请求被 Gin 封装在 Gin Context 里
- 然后经过 Logger 和 Recovery 中间件
  - logger 会在 stdout 输出请求的信息
  - recovery 会在 panic 的时候，恢复到可以接受请求的状态，然后给客户端一个 500 错误。
- 然后经过 GinContextToContextMiddleware 中间件
  - 请求会从 GinContext 放到 Context 里面
  - 为了 Gql resolver 可以访问请求
- 请求根据预先设定的路由，到达不同的路径。
- 最主要的请求到达 /query 路径。
- 然后会被 gqlgen 解析
