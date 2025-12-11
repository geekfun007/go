package client

import (
	"log"
	"sync"

	"github.com/cloudwego/kitex/client"
	etcd "github.com/kitex-contrib/registry-etcd"

	"github.com/example/microservices-project/api-gateway/kitex_gen/user/userservice"
)

var (
	userClient userservice.Client
	once       sync.Once
)

func InitUserClient(etcdEndpoints []string) {
	once.Do(func() {
		// 创建 Etcd 解析器
		r, err := etcd.NewEtcdResolver(etcdEndpoints)
		if err != nil {
			log.Fatal("Failed to create etcd resolver:", err)
		}

		// 创建 Kitex 客户端
		c, err := userservice.NewClient(
			"user-service",
			client.WithResolver(r),
		)
		if err != nil {
			log.Fatal("Failed to create user client:", err)
		}

		userClient = c
		log.Println("User client initialized successfully")
	})
}

func GetUserClient() userservice.Client {
	if userClient == nil {
		log.Fatal("User client not initialized")
	}
	return userClient
}
