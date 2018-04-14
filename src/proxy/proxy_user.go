package proxy

import (
	"io"
	"log"
	"net"
	"time"

	hz "github.com/dustin/go-humanize"
	"github.com/xsank/EasyProxy/src/structure"
)

func (proxy *EasyProxy) transfer2(local net.Conn, remote string) {
	user, buf := ParseHeader(local)
	if !isValidUser(user) {
		remote = FHOST
		user = ""
		//log.Println("change remote to", remote)
	} else {
		remote = changeRemote(user, remote)
		//log.Println("change remote to", remote)
	}

	remoteConn, err := net.DialTimeout("tcp", remote, DefaultTimeoutTime*time.Second)
	if err != nil {
		local.Close()
		proxy.Clean(remote)
		log.Printf("connect backend error:%s\n", err)
		return
	}

	if _, err := remoteConn.Write(buf); err != nil {
		local.Close()
		proxy.Clean(remote)
		log.Printf("connect backend error:%s\n", err)
		return
	}
	sync := make(chan int, 1)
	channel := structure.Channel{SrcConn: local, DstConn: remoteConn}
	go proxy.putChannel(&channel)
	if user != "" {
		go proxy.safeCopy2(local, remoteConn, sync, user)
	} else {
		go proxy.safeCopy(local, remoteConn, sync)
	}
	go proxy.safeCopy(remoteConn, local, sync)
	go proxy.closeChannel(&channel, sync)
}

func (proxy *EasyProxy) safeCopy2(from net.Conn, to net.Conn, sync chan int, user string) {
	n, _ := io.Copy(from, to)
	Users[user].Traf += uint64(n)
	Users[user].TrafficLeft -= uint64(n)
	Users[user].CurrentRemote = from.RemoteAddr().String()

	log.Println(user, ":", hz.Bytes(uint64(Users[user].Traf)), "|", hz.Bytes(uint64(Users[user].TrafficLeft)), "via:", from.RemoteAddr())
	defer from.Close()
	sync <- 1
}
