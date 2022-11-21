package main

import (
	_ "aiyun_cloud_srv/boot"
	_ "aiyun_cloud_srv/router"
	"github.com/gogf/gf/frame/g"
	_ "github.com/lib/pq"
)

func main() {
	//peers := make([]net.UDPAddr, 0, 2)
	//data := make([]byte, 1024)
	//// Server
	//go gudp.NewServer("0.0.0.0:9999", func(conn *gudp.Conn) {
	//	g.Log().Info("UDP")
	//	defer conn.Close()
	//	for {
	//		n, RemoteAddr, err := conn.ReadFromUDP(data)
	//		if err != nil {
	//			g.Log().Errorf("err during read: %s", err.Error())
	//		}
	//		g.Log().Infof("%s:%s", RemoteAddr.String(), data[:n])
	//		peers = append(peers, *RemoteAddr)
	//
	//		if len(peers) == 2 {
	//			conn.WriteToUDP([]byte(peers[1].String()), &peers[0])
	//			conn.WriteToUDP([]byte(peers[0].String()), &peers[1])
	//			time.Sleep(time.Second * 8)
	//			g.Log().Info("UDP打洞，中转服务器退出不影响peers通信")
	//			return
	//		}
	//	}
	//}).Run()

	g.Server().Run()
}
