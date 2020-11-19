module github.com/imgk/shadow

go 1.15

require (
	github.com/eycorsican/go-tun2socks v1.16.11
	github.com/gorilla/websocket v1.4.2
	github.com/imgk/divert-go v0.0.0-20201117053927-7aaa9bed883f
	github.com/miekg/dns v1.1.35
	github.com/oschwald/maxminddb-golang v1.7.0
	github.com/xtaci/smux v1.5.14
	go.uber.org/multierr v1.6.0
	go.uber.org/zap v1.16.0
	golang.org/x/crypto v0.0.0-20201117144127-c1f2f97bffc9
	golang.org/x/net v0.0.0-20201110031124-69a78807bb2b
	golang.org/x/sys v0.0.0-20201118182958-a01c418693c7
	golang.org/x/time v0.0.0-20200630173020-3af7569d3a1e
	golang.zx2c4.com/wireguard v0.0.20201118
	golang.zx2c4.com/wireguard/windows v0.2.2
	gvisor.dev/gvisor v0.0.0-20201119000631-7158095d687d
)

replace gvisor.dev/gvisor => gvisor.dev/gvisor v0.0.0-20201119071348-817a3a5fa787
