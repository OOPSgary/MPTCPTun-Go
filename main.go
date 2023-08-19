package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"strings"
	"sync"
)

//because im poor at coding, so =(

var (
	dialer = &net.Dialer{}
	lc     = &net.ListenConfig{}
)

func init() {
	// Set MultiPath into True
	// According to Golang Official Doc,this wont be automatically enabled in 1.21
	// If client or server cant support this(MPTCP) feature, the MPTCP will fall back to Normal TCP
	dialer.SetMultipathTCP(true)
	lc.SetMultipathTCP(true)
}
func main() {
	var querylist stringList
	help := flag.Bool("h", false, "Show this message")
	flag.Var(&querylist, "L", "Targets(:ListenPort/Server1:Server1P/S2:S2P")
	servermode := flag.Bool("s", false, "Servermode(Default: False)")
	log.Println("WARNING: Only Linux kernels >= v5.16 CAN USE MPTCP")
	flag.Parse()
	if querylist == nil {
		flag.PrintDefaults()
		return
	}
	if *help {
		flag.PrintDefaults()
		return
	}
	var wait *sync.WaitGroup = &sync.WaitGroup{}

	for _, t := range querylist {
		wait.Add(1)
		slice := strings.Split(t, "/")
		if len(slice) != 2 {
			log.Fatalf("Check your Inputs!:%v", t)
		}
		if *servermode {
			go server(slice[0], slice[1:], wait)
		} else {
			go client(slice[0], slice[1:], wait)
		}
	}
	wait.Wait()
}

func server(listenPort string, ServerAddress []string, wg *sync.WaitGroup) {
	defer wg.Done()
	ln, err := lc.Listen(context.Background(), "tcp", fmt.Sprintf("[::]%v", listenPort))
	if err != nil {
		log.Fatalf("listen tcp error:%v", err)
	}
	log.Printf("Listening on Address: %v\n", ln.Addr().String())
	for {
		c, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		tconn, _ := c.(*net.TCPConn)
		isMultipathTCP, err := tconn.MultipathTCP()
		if !isMultipathTCP {
			log.Printf("WARNING: MPTCP NOT ENABLED %v\n", err)
		}
		log.Printf("New Connection Connected: %v\n", c.RemoteAddr().String())
		go handlerServer(c, ServerAddress[rand.Intn(len(ServerAddress))])
	}
}
func client(listenPort string, ServerAddress []string, wg *sync.WaitGroup) {
	defer wg.Done()
	listentcp, err := lc.Listen(context.Background(), "tcp", fmt.Sprintf("[::]%v", listenPort))
	if err != nil {
		log.Fatalf("listen tcp error:%v", err)
	}
	for {
		c, err := listentcp.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		log.Printf("New Connection Connected: %v\n", c.RemoteAddr().String())
		go handlerClientTCP(c, ServerAddress[rand.Intn(len(ServerAddress))])
	}
}

// Inspired from https://github.com/ginuerzh/gost/blob/master/cmd/gost/route.go#L19
type stringList []string

func (l *stringList) String() string {
	return fmt.Sprintf("%s", *l)
}
func (l *stringList) Set(value string) error {
	*l = append(*l, value)
	return nil
}

func handlerServer(c net.Conn, target string) {
	defer c.Close()

	remote, err := dialer.Dial("tcp", target)
	if err != nil {
		log.Println(err)
		return
	}

	defer remote.Close()
	go func() {
		defer remote.Close()
		defer c.Close()
		io.Copy(remote, c)
	}()
	io.Copy(c, remote)
}

func handlerClientTCP(c net.Conn, target string) {
	remote, err := dialer.Dial("tcp", target)
	if err != nil {
		log.Println(err)
		return
	}
	isMultipathTCP, err := remote.(*net.TCPConn).MultipathTCP()
	if !isMultipathTCP {
		log.Printf("WARNING: MPTCP NOT ENABLED %v\n", err)
	}
	defer remote.Close()
	go func() {
		defer remote.Close()
		defer c.Close()
		io.Copy(remote, c)
	}()
	io.Copy(c, remote)
}
