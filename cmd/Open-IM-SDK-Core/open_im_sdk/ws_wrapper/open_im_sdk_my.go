/*
** description("").
** copyright('open-im,www.open-im.io').
** author("fg,Gordon@tuoyun.net").
** time(2021/9/8 14:35).
 */
package main

import (
	"flag"
	"fmt"
	"open_im_sdk/open_im_sdk"
	"open_im_sdk/open_im_sdk/ws_wrapper/utils"
	"open_im_sdk/open_im_sdk/ws_wrapper/ws_local_server"
	"runtime"
	"sync"
)

func main() {
	var sdkWsPort, openIMApiPort, openIMWsPort *int
	var openIMWsAddress, openIMApiAddress *string
	APIADDR := "http://47.112.160.66:10000"
	WSADDR := "ws://47.112.160.66:17778"

	sysType := runtime.GOOS
	switch sysType {
	case "darwin":
		ws_local_server.InitServer(&open_im_sdk.IMConfig{IpApiAddr: *openIMApiAddress,
			IpWsAddr: *openIMWsAddress, Platform: utils.OSXPlatformID, DbDir: "./"})
	case "linux":

		//sdkDBDir:= flag.String("sdk_db_dir","","openIMSDK initialization path")
		ws_local_server.InitServer(&open_im_sdk.IMConfig{IpApiAddr: "http://" + utils.ServerIP + ":" + utils.IntToString(*openIMApiPort),
			IpWsAddr: "ws://" + utils.ServerIP + ":" + utils.IntToString(*openIMWsPort), Platform: utils.WebPlatformID, DbDir: "../db/sdk/"})

	case "windows":
		sdkWsPort = flag.Int("sdk_ws_port", 7799, "openIM ws listening port")
		flag.Parse()
		ws_local_server.InitServer(&open_im_sdk.IMConfig{IpApiAddr: APIADDR,
			IpWsAddr: WSADDR, Platform: utils.WebPlatformID, DbDir: "./"})
	default:
		fmt.Println("this os not support", sysType)

	}
	var wg sync.WaitGroup
	wg.Add(1)
	fmt.Println("ws server is starting")
	ws_local_server.WS.OnInit(*sdkWsPort)
	ws_local_server.WS.Run()
	wg.Wait()

}
