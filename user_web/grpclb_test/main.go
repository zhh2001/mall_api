package main

import (
	"context"
	"fmt"
	"log"

	_ "github.com/mbobakov/grpc-consul-resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"mall_api/user_web/proto"
)

func main() {
	conn, err := grpc.NewClient(
		"consul://10.120.21.77:8500/user_srv?wait=14s&tag=srv",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(conn)

	userSrvClient := proto.NewUserClient(conn)
	for i := 0; i < 10; i++ {
		rsp, err := userSrvClient.GetUserList(context.Background(), &proto.PageInfo{
			Pn:    1,
			PSize: 2,
		})
		if err != nil {
			log.Fatal(err)
		}
		for index, data := range rsp.Data {
			fmt.Println(index, data)
		}
	}
}
