/*
 * Copyright (c) 2022, Alibaba Group;
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"context"
	"fmt"
	"time"

	"github.com/alibaba/ioc-golang"
	"github.com/alibaba/ioc-golang/autowire"
	"github.com/alibaba/ioc-golang/extension/autowire/grpc"

	googleGRPC "google.golang.org/grpc"

	"github.com/alibaba/ioc-golang/example/debug/api"
	"github.com/alibaba/ioc-golang/example/debug/cmd/service1"
	"github.com/alibaba/ioc-golang/example/debug/cmd/service2"
	"github.com/alibaba/ioc-golang/example/debug/cmd/struct1"
)

func init() {
	// register grpc client
	grpc.RegisterStructDescriptor(&autowire.StructDescriptor{
		Factory: func() interface{} {
			return new(api.HelloServiceClient)
		},
		ParamFactory: func() interface{} {
			return &googleGRPC.ClientConn{}
		},
		ConstructFunc: func(impl interface{}, param interface{}) (interface{}, error) {
			conn := param.(*googleGRPC.ClientConn)
			fmt.Println("create conn target ", conn.Target())
			return api.NewHelloServiceClient(conn), nil
		},
	})
}

// +ioc:autowire=true
// +ioc:autowire:type=singleton

type App struct {
	HelloServiceClient api.HelloServiceClient `grpc:"hello-service"`

	ExampleService1Impl1 service1.Service1 `singleton:"github.com/alibaba/ioc-golang/example/debug/cmd/service1.Impl1"`

	ExampleService2Impl1 service2.Service2 `singleton:"github.com/alibaba/ioc-golang/example/debug/cmd/service2.Impl1"`
	ExampleService2Impl2 service2.Service2 `singleton:"github.com/alibaba/ioc-golang/example/debug/cmd/service2.Impl2"`

	ExampleStruct1 *struct1.Struct1 `singleton:""`
}

func (a *App) Run() {
	name := "laurence"
	rsp, err := a.HelloServiceClient.SayHello(context.Background(), &api.HelloRequest{
		Name: name,
	}, googleGRPC.MaxCallRecvMsgSize(1024), googleGRPC.MaxCallRecvMsgSize(1024))
	if err != nil {
		panic(err)
	}

	fmt.Println("App call grpc get: " + rsp.Reply)

	fmt.Println("ExampleService1Impl1 call grpc get :" + a.ExampleService1Impl1.Hello(name))

	fmt.Println("ExampleService2Impl1 call grpc get :" + a.ExampleService2Impl1.Hello(name))

	fmt.Println("ExampleService2Impl2 call grpc get :" + a.ExampleService2Impl2.Hello(name))

	fmt.Println("ExampleStruct1 call grpc get :" + a.ExampleStruct1.Hello(name))
}

func main() {
	if err := ioc.Load(); err != nil {
		panic(err)
	}
	app, err := GetApp()
	if err != nil {
		panic(err)
	}

	for {
		time.Sleep(time.Second * 5)
		app.Run()
	}
}
