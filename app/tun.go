// +build linux darwin windows,!divert

package app

import (
	"fmt"
	"net"
	"net/http"
	"net/http/pprof"

	"github.com/imgk/shadow/pkg/handler/recorder"
	"github.com/imgk/shadow/pkg/netstack"
	"github.com/imgk/shadow/pkg/proxy"
	"github.com/imgk/shadow/pkg/resolver"
	"github.com/imgk/shadow/pkg/tun"
	"github.com/imgk/shadow/proto"
)

// RunWithDevice is ...
func (app *App) RunWithDevice(dev *tun.Device) (err error) {
	// new dns resolver
	resolver, err := resolver.NewResolver(app.Conf.NameServer)
	if err != nil {
		return fmt.Errorf("dns server error: %w", err)
	}
	net.DefaultResolver = &net.Resolver{
		PreferGo: true,
		Dial:     resolver.DialContext,
	}

	// new connection handler
	handler, err := proto.NewHandler(app.Conf.Server, app.Timeout)
	if err != nil {
		return fmt.Errorf("protocol error: %w", err)
	}
	handler = recorder.NewHandler(handler)
	app.attachCloser(handler)
	defer func() {
		if err != nil {
			for _, closer := range app.closers {
				closer.Close()
			}
		}
	}()

	router := http.NewServeMux()
	router.HandleFunc("/debug/pprof/", pprof.Index)
	router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	router.HandleFunc("/debug/pprof/profile", pprof.Profile)
	router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	router.HandleFunc("/debug/pprof/trace", pprof.Trace)
	router.Handle("/admin/conns", handler.(*recorder.Handler))
	router.Handle("/admin/proxy.pac", NewPACForSocks5())

	// new tun device
	name := "utun"
	if tunName := app.Conf.TunName; tunName != "" {
		name = tunName
	}
	createDevice := false
	if dev == nil {
		createDevice = true
		dev, err = tun.NewDeviceWithMTU(name, (2<<10)-4 /*MTU for Tun*/)
		if err != nil {
			return fmt.Errorf("tun device from name error: %w", err)
		}
	}
	app.attachCloser(dev)
	// set tun address
	for _, address := range app.Conf.TunAddr {
		err := dev.SetInterfaceAddress(address)
		if err != nil {
			return err
		}
	}
	if createDevice {
		if err := dev.Activate(); err != nil {
			return fmt.Errorf("turn up tun device error: %w", err)
		}
	}

	// new fake ip tree
	tree, err := NewDomainTree(app.Conf)
	if err != nil {
		return
	}
	// new netstack
	stack := netstack.NewStack(handler, resolver, tree, !app.Conf.DomainRules.DisableHijack /* ture for hijacking queries */)
	err = stack.Start(dev, app.Logger, (2<<10)-4 /*MTU for Tun*/)
	if err != nil {
		return
	}
	app.attachCloser(stack)

	// enable socks5/http proxy
	if addr := app.Conf.ProxyServer; addr != "" {
		ln, err := net.Listen("tcp", addr)
		if err != nil {
			return err
		}

		server := proxy.NewServer(ln, app.Logger, handler, tree, router)
		app.attachCloser(server)
		go server.Serve()
	}

	// add route table entry
	if err := dev.AddRouteEntry(app.Conf.IPCIDRRules.Proxy); err != nil {
		return fmt.Errorf("add route entry error: %w", err)
	}

	return nil
}

// Run is ...
func (app *App) Run() error {
	return app.RunWithDevice(nil)
}
