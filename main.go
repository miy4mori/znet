package main

import (
	"bytes"
	"log"
	"runtime"
	"znet/epoll"
	"znet/listener"
	"znet/loop"
)

func main() {
	p := loop.NewEventLoopPool(runtime.NumCPU())

	e, err := epoll.NewEpoll()
	if err != nil {
		return
	}
	defer func() { _ = e.Close() }()

	l := listener.NewListener()
	defer func() { _ = l.Close() }()

	_ = l.Listen(9999)

	p.Run()

	go func() {
		for {
			connections, err := e.Wait()
			if err != nil {
				log.Printf("Wait Error %v\n", err)
			}
			for _, connection := range connections {
				buffer, err := connection.Read()
				if err != nil {
					log.Printf("Read Error %v\n", err)
					continue
				}
				t := loop.NewTaskWithProcess(buffer, func(buffer *bytes.Buffer) {
					s, err := buffer.ReadString('\n')
					if err != nil {
						log.Printf("Read String Error")
						return
					}
					log.Printf("Read Data %s", s)
				})
				p.Execute(t)
			}
		}
	}()

	for {
		c, _ := l.Accept()
		_ = e.Add(c)
	}
}
