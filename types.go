package main

type (
	// Client- настройки клиента
	Client struct {
		Servers      []string
		TunnelPort   string
		LocalPort    string
		LocalAdress  string
		TargetUrl    string
		Dom          string
		CrtFile      string
		TokenCookieA string
		TokenCookieB string
		TokenCookieC string
		UserAgent    string
		WsObf        bool
		TlsVerify    bool
	}

	Server struct {
		Port         string
		Target       string
		Dir          string
		TokenCookieA string
		TokenCookieB string
		TokenCookieC string
		HeaderServer string
		WsObf        bool
		OnlyWs       bool
		CrtFile      string
		KeyFile      string
	}

	// Config - настройки
	Config struct {
		Client Client `json:"Client"`
		Server Server `json:"Server"`
	}
)
