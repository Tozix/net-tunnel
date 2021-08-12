package main

import (
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net-tunnel/fakehttp"
	"os"
	"runtime"
)

func handleClient(p1 net.Conn, servers []string, port string) {
	defer p1.Close()

	p2, err := cl.Dial()
	if err != nil {
		Vlogln(2, cl.Host, "- Сервер не доступен")
		for _, server := range servers {
			Vlogln(2, "Конектится к: ", server+":"+port)
			cl.Host = server + ":" + port
			p2, err = cl.Dial()
			if err != nil {
				Vlogln(2, "Dial err2:", err)
				continue
			} else {
				break
			}
		}
		return
	}

	defer p2.Close()
	cp(p1, p2)
	Vlogln(2, "close", p1.RemoteAddr())
}

func client(cfg Client) {
	log.Println("Стартуем клиент.")

	runtime.GOMAXPROCS(runtime.NumCPU() + 2)
	flag.Parse()

	if cfg.TokenCookieA == cfg.TokenCookieB {
		Vlogln(2, "Error: token cookie cannot bee same!")
		os.Exit(1)
	}

	lis, err := net.Listen("tcp", cfg.LocalAdress+":"+cfg.LocalPort)
	if err != nil {
		Vlogln(2, "Error listening:", err.Error())
		os.Exit(1)
	}
	defer lis.Close()

	Vlogln(2, "Слушаем на:", lis.Addr())
	Vlogln(2, "Домен для коннекта:", cfg.Dom)
	Vlogln(2, "token cookie A:", cfg.TokenCookieA)
	Vlogln(2, "token cookie B:", cfg.TokenCookieB)
	Vlogln(2, "token cookie C:", cfg.TokenCookieC)
	Vlogln(2, "use ws:", cfg.WsObf)

	if cfg.CrtFile != "" {
		caCert, err := ioutil.ReadFile(cfg.CrtFile)
		if err != nil {
			Vlogln(2, "Reading certificate error:", err)
			os.Exit(1)
		}
		cl = fakehttp.NewTLSClient(cfg.Servers[0]+":"+cfg.TunnelPort, caCert, cfg.TlsVerify)
	} else {
		Vlogln(2, "Тунель готов к подключению")
		cl = fakehttp.NewClient(cfg.Servers[0] + ":" + cfg.TunnelPort)
	}
	cl.TokenCookieA = cfg.TokenCookieA
	cl.TokenCookieB = cfg.TokenCookieB
	cl.TokenCookieC = cfg.TokenCookieC
	cl.UseWs = cfg.WsObf
	cl.UserAgent = cfg.UserAgent
	cl.CDom = cfg.Dom
	cl.Url = cfg.TargetUrl

	copyBuf.New = func() interface{} {
		return make([]byte, 4096)
	}

	for {
		if conn, err := lis.Accept(); err == nil {
			Vlogln(2, "remote address:", conn.RemoteAddr())

			go handleClient(conn, cfg.Servers, cfg.TunnelPort)
		} else {
			Vlogf(2, "%+v", err)
		}
	}

}

func cp(p1, p2 io.ReadWriteCloser) {
	//	Vlogln(2, "stream opened")
	//	defer Vlogln(2, "stream closed")
	//	defer p1.Close()
	//	defer p2.Close()

	// start tunnel
	p1die := make(chan struct{})
	go func() {
		buf := copyBuf.Get().([]byte)
		io.CopyBuffer(p1, p2, buf)
		close(p1die)
		copyBuf.Put(buf)
	}()

	p2die := make(chan struct{})
	go func() {
		buf := copyBuf.Get().([]byte)
		io.CopyBuffer(p2, p1, buf)
		close(p2die)
		copyBuf.Put(buf)
	}()

	// wait for tunnel termination
	select {
	case <-p1die:
	case <-p2die:
	}
}
