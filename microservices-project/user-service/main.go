package main

import (
	"fmt"
	"log"
	"net"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"

	"github.com/example/microservices-project/user-service/biz/handler"
	"github.com/example/microservices-project/user-service/biz/model"
	"github.com/example/microservices-project/user-service/biz/repository"
	"github.com/example/microservices-project/user-service/biz/service"
	"github.com/example/microservices-project/user-service/conf"
	"github.com/example/microservices-project/user-service/kitex_gen/user/userservice"
	"github.com/example/microservices-project/user-service/pkg/db"
)

func main() {
	// 加载配置
	cfg, err := conf.LoadConfig("conf/config.yaml")
	if err != nil {
		log.Printf("Failed to load config: %v, using default config", err)
		cfg = conf.GetConfig()
	}

	// 初始化数据库
	database, err := db.InitMySQL(&cfg.Database)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// 自动迁移
	if err := database.AutoMigrate(&model.User{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	klog.Info("Database migrated successfully")

	// 初始化依赖
	userRepo := repository.NewUserRepository(database)
	userSvc := service.NewUserService(userRepo)
	userHandler := handler.NewUserService(userSvc)

	// 创建 Etcd 注册中心
	r, err := etcd.NewEtcdRegistry(cfg.Etcd.Endpoints)
	if err != nil {
		log.Fatal("Failed to create etcd registry:", err)
	}
	klog.Infof("Etcd registry created with endpoints: %v", cfg.Etcd.Endpoints)

	// 服务地址
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", cfg.Server.Address, cfg.Server.Port))
	if err != nil {
		log.Fatal("Failed to resolve address:", err)
	}

	// 创建 Kitex 服务器
	svr := userservice.NewServer(
		userHandler,
		server.WithServiceAddr(addr),
		server.WithRegistry(r),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
			ServiceName: cfg.Server.ServiceName,
		}),
	)

	klog.Infof("User service starting on %s:%d", cfg.Server.Address, cfg.Server.Port)
	klog.Infof("Service name: %s", cfg.Server.ServiceName)

	err = svr.Run()
	if err != nil {
		log.Fatal(err)
	}
}
