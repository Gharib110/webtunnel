package main

import (
	"errors"
	"fmt"
	"gitlab.torproject.org/tpo/anti-censorship/pluggable-transports/webtunnel/common/syntheticIP"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"gitlab.torproject.org/tpo/anti-censorship/pluggable-transports/webtunnel"

	pt "gitlab.torproject.org/tpo/anti-censorship/pluggable-transports/goptlib"
)

const ptMethodName = "webtunnel"

func main() {
	var err error
	var ptInfo pt.ServerInfo

	ptInfo, err = pt.ServerSetup(nil)
	if err != nil {
		log.Fatalf("error in setup: %s", err)
	}
	pt.ReportVersion("webtunnel", webtunnel.Version)

	listeners := make([]net.Listener, 0)
	for _, bindaddr := range ptInfo.Bindaddrs {
		if bindaddr.MethodName != ptMethodName {
			pt.SmethodError(bindaddr.MethodName, "no such method")
			continue
		}

		args := bindaddr.Options
		config := &ServerConfig{}
		ln, err := net.ListenTCP("tcp", bindaddr.Addr)
		if err != nil {
			log.Fatalf("error in setup: %s", err)
		}
		defer ln.Close()
		go acceptLoop(ln, config, &ptInfo)

		if args == nil {
			args = pt.Args{}
		}

		args.Add("ver", webtunnel.Version)

		urlValue, ok := args.Get("url")
		if !ok {
			pt.SmethodError(bindaddr.MethodName, "missing url parameter")
			continue
		}

		_, cidr, err := net.ParseCIDR("2001:DB8::/32")
		if err != nil {
			pt.SmethodError(bindaddr.MethodName, fmt.Sprintf("error in ParseCIDR: %s", err))
			continue
		}

		generatedAddress, err := syntheticIP.GenerateSyntheticIPAddress("WEBTUNNEL+"+urlValue, *cidr)
		if err != nil {
			pt.SmethodError(bindaddr.MethodName, fmt.Sprintf("error in GenerateSyntheticIPAddress: %s", err))
			continue
		}

		generatedAddr := &net.TCPAddr{IP: generatedAddress, Port: 443}
		pt.SmethodArgs(bindaddr.MethodName, generatedAddr, args)
		listeners = append(listeners, ln)
	}

	if len(listeners) == 0 {
		pt.SmethodError(ptMethodName, "no valid listener configured")
		return
	}

	pt.SmethodsDone()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM)

	if os.Getenv("TOR_PT_EXIT_ON_STDIN_CLOSE") == "1" {
		// This environment variable means we should treat EOF on stdin
		// just like SIGTERM: https://bugs.torproject.org/15435.
		go func() {
			if _, err := io.Copy(ioutil.Discard, os.Stdin); err != nil {
				log.Printf("error copying os.Stdin to ioutil.Discard: %v", err)
			}
			log.Printf("synthesizing SIGTERM because of stdin close")
			sigChan <- syscall.SIGTERM
		}()
	}

	// Wait for a signal.
	sig := <-sigChan

	// Signal received, shut down.
	log.Printf("caught signal %q, exiting", sig)
	for _, ln := range listeners {
		ln.Close()
	}

}

func acceptLoop(ln net.Listener, config *ServerConfig, ptInfo *pt.ServerInfo) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			if err, ok := err.(net.Error); ok && err.Temporary() {
				continue
			}
			log.Printf("webtunnel accept error: %s", err)
			break
		}
		go func() {
			defer conn.Close()

			transport, err := NewWebTunnelServerTransport(config)
			conn, err := transport.Accept(conn)
			if err != nil {
				log.Printf("handleConn: %v", err)
				return
			}
			handleConn(conn, ptInfo)
		}()
	}
}

// proxy copies data bidirectionally from one connection to another.
func proxy(local *net.TCPConn, conn net.Conn) {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		if _, err := io.Copy(conn, local); err != nil && !errors.Is(err, io.ErrClosedPipe) {
			log.Printf("error copying ORPort to WebSocket %v", err)
		}
		local.CloseRead()
		conn.Close()
		wg.Done()
	}()
	go func() {
		if _, err := io.Copy(local, conn); err != nil && !errors.Is(err, io.ErrClosedPipe) {
			log.Printf("error copying WebSocket to ORPort %v", err)
		}
		local.CloseWrite()
		conn.Close()
		wg.Done()
	}()

	wg.Wait()
}

// handleConn bidirectionally connects a client webtunnel connection with an ORPort.
func handleConn(conn net.Conn, ptInfo *pt.ServerInfo) error {
	addr := conn.RemoteAddr().String()
	or, err := pt.DialOr(ptInfo, addr, ptMethodName)
	if err != nil {
		return fmt.Errorf("failed to connect to ORPort: %s", err)
	}
	defer or.Close()
	proxy(or, conn)
	return nil
}
