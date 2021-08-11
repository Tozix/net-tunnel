package main

import (
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"runtime"
	"sync"

	"./fakehttp"
)
//nts 95.174.206.116
//tomtel 79.136.197.163
//rt 90.188.116.11
//boost 176.120.28.123
var verbosity int = 2

var copyBuf sync.Pool
//var servers =: []string{79.136.197.163,176.120.28.123,95.174.206.116,90.188.116.11} 

var servers = []string{"90.188.116.11"}
var port = flag.String("p", "127.0.0.1:5555", "bind port")
var target = flag.String("t", "90.188.116.11:80", "http server address & port")
var targetUrl = flag.String("url", "/tunnel", "http url to send")

var crtFile = flag.String("crt", "", "PEM encoded certificate file")

var tokenCookieA = flag.String("ca", "cna", "token cookie name A")
var tokenCookieB = flag.String("cb", "_tb_token_", "token cookie name B")
var tokenCookieC = flag.String("cc", "_cna", "token cookie name C")

var userAgent = flag.String("ua", "BradburyLab (samsung_dreamqltevl; SM-G930L; Android; 5.1.1) 7.15.0 (22064)", "User-Agent (default: QQ)")
var cDom = flag.String("dom", "nsk.mts.ru", "Хост для коннекта")

var wsObf = flag.Bool("usews", false, "fake as websocket")
var tlsVerify = flag.Bool("k", true, "InsecureSkipVerify")

var cl *fakehttp.Client

func handleClient(p1 net.Conn) {
	defer p1.Close()
  
	p2, err := cl.Dial()
    if err != nil {
        Vlogln(2,cl.Host, "- Сервер не доступен")
    for _, server := range servers {
        Vlogln(2,"Конектится к: ", server)
        cl.Host=server+":80"
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

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() + 2)
	flag.Parse()

	if *tokenCookieA == *tokenCookieB {
		Vlogln(2, "Error: token cookie cannot bee same!")
		os.Exit(1)
	}

	lis, err := net.Listen("tcp", *port)
	if err != nil {
		Vlogln(2, "Error listening:", err.Error())
		os.Exit(1)
	}
	defer lis.Close()

	Vlogln(2, "listening on:", lis.Addr())
	Vlogln(2, "target:", *target)
	Vlogln(2, "Домен для коннекта:", *cDom)
	Vlogln(2, "token cookie A:", *tokenCookieA)
	Vlogln(2, "token cookie B:", *tokenCookieB)
	Vlogln(2, "token cookie C:", *tokenCookieC)
	Vlogln(2, "use ws:", *wsObf)

	if *crtFile != "" {
		caCert, err := ioutil.ReadFile(*crtFile)
		if err != nil {
			Vlogln(2, "Reading certificate error:", err)
			os.Exit(1)
		}
		cl = fakehttp.NewTLSClient(*target, caCert, *tlsVerify)
	} else {
		Vlogln(2, "Тунель готов к подключению")
		cl = fakehttp.NewClient(*target)
	}
	cl.TokenCookieA = *tokenCookieA
	cl.TokenCookieB = *tokenCookieB
	cl.TokenCookieC = *tokenCookieC
	cl.UseWs = *wsObf
	cl.UserAgent = *userAgent
	cl.CDom = *cDom
	cl.Url = *targetUrl

	copyBuf.New = func() interface{} {
		return make([]byte, 4096)
	}

	for {
		if conn, err := lis.Accept(); err == nil {
			Vlogln(2, "remote address:", conn.RemoteAddr())

			go handleClient(conn)
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
func Vlogf(level int, format string, v ...interface{}) {
	if level <= verbosity {
		log.Printf(format, v...)
	}
}
func Vlog(level int, v ...interface{}) {
	if level <= verbosity {
		log.Print(v...)
	}
}
func Vlogln(level int, v ...interface{}) {
	if level <= verbosity {
		log.Println(v...)
	}
}
