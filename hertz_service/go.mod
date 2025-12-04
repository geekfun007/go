module github.com/example/hertz-kitex-demo

go 1.23.0

toolchain go1.24.11

replace github.com/example/hertz-kitex-demo => ../

require (
	github.com/cloudwego/hertz v0.10.3
	github.com/cloudwego/kitex v0.15.2
	github.com/kitex-contrib/registry-etcd v0.3.0
)
