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

package interceptor

import (
	"net"

	"github.com/fatih/color"
	"google.golang.org/grpc"

	"github.com/alibaba/ioc-golang/common"

	"github.com/alibaba/ioc-golang/debug/api/ioc_golang/boot"
	debugCommon "github.com/alibaba/ioc-golang/debug/common"
)

func Start(port string, allInterfaceMetadataMap map[string]*debugCommon.StructMetadata) error {
	grpcServer := grpc.NewServer()
	grpcServer.RegisterService(&boot.DebugService_ServiceDesc, &DebugServerImpl{
		editInterceptor:         GetEditInterceptor(),
		watchInterceptor:        GetWatchInterceptor(),
		allInterfaceMetadataMap: allInterfaceMetadataMap,
	})

	lst, err := common.GetTCPListener(port)
	if err != nil {
		return err
	}

	go func() {
		color.Blue("[Debug] Debug server listening at :%s", lst.Addr().(*net.TCPAddr).Port)
		if err := grpcServer.Serve(lst); err != nil {
			color.Red("[Debug] Debug server run with error = ", err)
			return
		}
	}()
	return nil
}
