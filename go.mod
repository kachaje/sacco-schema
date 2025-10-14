module github.com/kachaje/sacco-schema

go 1.24.3

require gopkg.in/yaml.v3 v3.0.1

require (
	github.com/google/uuid v1.6.0
	github.com/gorilla/mux v1.8.1
	github.com/gorilla/websocket v1.5.3
	github.com/joselitofilho/drawio-parser-go v0.3.2
	github.com/kachaje/firestore-webrtc v0.0.121
	github.com/rogpeppe/go-internal v1.14.1
	golang.org/x/crypto v0.42.0
	golang.org/x/exp v0.0.0-20250620022241-b7579e27df2b
	golang.org/x/net v0.44.0
	golang.org/x/text v0.29.0
	modernc.org/sqlite v1.39.0
)

require (
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/ncruces/go-strftime v0.1.9 // indirect
	github.com/pion/datachannel v1.5.10 // indirect
	github.com/pion/dtls/v3 v3.0.7 // indirect
	github.com/pion/ice/v4 v4.0.10 // indirect
	github.com/pion/interceptor v0.1.40 // indirect
	github.com/pion/logging v0.2.4 // indirect
	github.com/pion/mdns/v2 v2.0.7 // indirect
	github.com/pion/randutil v0.1.0 // indirect
	github.com/pion/rtcp v1.2.15 // indirect
	github.com/pion/rtp v1.8.21 // indirect
	github.com/pion/sctp v1.8.39 // indirect
	github.com/pion/sdp/v3 v3.0.15 // indirect
	github.com/pion/srtp/v3 v3.0.7 // indirect
	github.com/pion/stun/v3 v3.0.0 // indirect
	github.com/pion/transport/v3 v3.0.7 // indirect
	github.com/pion/turn/v4 v4.1.1 // indirect
	github.com/pion/webrtc/v4 v4.1.4 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	github.com/wlynxg/anet v0.0.5 // indirect
	golang.org/x/sys v0.36.0 // indirect
	golang.org/x/tools v0.36.0 // indirect
	modernc.org/libc v1.66.3 // indirect
	modernc.org/mathutil v1.7.1 // indirect
	modernc.org/memory v1.11.0 // indirect
)

replace github.com/kachaje/firestore-webrtc => ../firestore-webrtc
