package main

import (
	"context"
	"fmt"
	cache "go-pattern/internal/cache/multilevel"
	"go-pattern/internal/config"
	"go-pattern/internal/initializer"
	"go-pattern/internal/model"
	repoFactory "go-pattern/internal/repo/factory"
	orderService "go-pattern/internal/service/order"
	productService "go-pattern/internal/service/product"
	"go-pattern/internal/table"
	"log"
	"time"
)

func main() {
	// 初始化数据库
	configs, err := config.InitConfig()
	if err != nil {
		log.Fatalf("InitConfig(): 初始化配置失败: %v", err)
	}
	// 连接数据库
	gormDB, err := initializer.GormDB(&configs.Database)
	if err != nil {
		log.Fatalf("GormDB: 数据库连接失败: %v", err)
	}
	// 创建所有表
	err = table.NewAllTables(gormDB)
	if err != nil {
		panic(err)
	}

	repoFactory := repoFactory.NewRepoFactory(gormDB)

	//userService := userService.NewUserService(repoFactory)
	orderService := orderService.NewOrderService(repoFactory)
	productService := productService.NewProductService(repoFactory)

	redis, err := initializer.Redis(&configs.Redis)
	if err != nil {
		log.Fatalf("Redis: 创建Redis失败: %v", err)
	}

	cacheFactory := cache.NewMultiLevelCacheFactory(redis)
	orderCache := cacheFactory.Order(&configs.LocalCache, 15*time.Second)
	productCache := cacheFactory.Product(&configs.LocalCache, 15*time.Second)

	err = productService.CreateProduct(context.Background(), &model.Product{
		Name:        "testproduct",
		Description: "testproduct",
		Price:       100,
		Quantity:    1000,
	})
	if err != nil {
		log.Fatalf("create product failed: %v", err)
	}

	productPointer, err := productService.GetProduct(context.Background(), 1)
	if err != nil {
		log.Fatalf("get product failed: %v", err)
	}
	log.Printf("product: %v", productPointer)

	productCache.SetWithDefaultTTL(context.Background(), fmt.Sprintf("product:%d", productPointer.ID), *productPointer)

	product, err := productCache.Get(context.Background(), fmt.Sprintf("product:%d", productPointer.ID))
	if err != nil {
		log.Fatalf("get product from cache failed: %v", err)
	}
	log.Printf("product from cache: %v", product)

	err = productService.ReduceQuantity(context.Background(), productPointer.ID, 100)
	if err != nil {
		log.Fatalf("reduce product quantity failed: %v", err)
	}
	productPointer, err = productService.GetProduct(context.Background(), productPointer.ID)
	if err != nil {
		log.Fatalf("get product failed: %v", err)
	}
	log.Printf("product after reduce quantity: %v", productPointer)

	productCache.Del(context.Background(), fmt.Sprintf("product:%d", productPointer.ID))

	product, err = productCache.Get(context.Background(), fmt.Sprintf("product:%d", productPointer.ID))
	if err != nil {
		log.Fatalf("get product from cache failed: %v", err)
	}
	log.Printf("product from cache after del: %v", product)

	orders, err := orderService.GetOrdersByUserID(context.Background(), 2)
	if err != nil {
		log.Fatalf("get orders failed: %v", err)
	}
	log.Printf("orders: %v", orders)
	for _, order := range orders {
		log.Printf("order: %v", order)
	}

	for _, order := range orders {
		orderCache.SetWithDefaultTTL(context.Background(), fmt.Sprintf("order:%d", order.ID), *order)
	}
	start := time.Now()
	for range 10000 {
		for _, order := range orders {
			_, err := orderCache.Get(context.Background(), fmt.Sprintf("order:%d", order.ID))
			if err != nil {
				log.Fatalf("get order from cache failed: %v", err)
			}
		}
	}
	log.Printf("cost: %v", time.Since(start))

	start = time.Now()
	for range 10000 {
		for _, order := range orders {
			_, err := orderService.GetOrder(context.Background(), order.ID)
			if err != nil {
				log.Fatalf("get order from db failed: %v", err)
			}
		}
	}
	log.Printf("cost: %v", time.Since(start))

}
