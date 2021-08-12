package main

import (
	"flag"
	"log"
	"net-tunnel/fakehttp"
	"runtime"
	"sync"
)

var verbosity int = 2
var copyBuf sync.Pool
var mode = flag.String("mode", "", "Режим работы: server или client")
var cl *fakehttp.Client

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() + 2)
	flag.Parse()
	// получение настроек /*

	config, err := GetConfig("config.json")
	if err != nil {
		log.Fatalln("Невозможно загрузить файл конфигурации:", err.Error())
	}
	log.Println("Запускаем." + config.Client.TokenCookieC)

	if *mode == "client" {
		client(config.Client)
	}
	if *mode == "server" {
		server(config.Server)
	}
	if *mode == "" {
		log.Fatalln("Не выбран режим работы туннеля")
	}

	log.Println("Запускаем.")

}
