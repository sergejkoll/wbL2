package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

/*
=== Утилита telnet ===

Реализовать примитивный telnet клиент:
Примеры вызовов:
go-telnet --timeout=10s host port go-telnet mysite.ru 8080 go-telnet --timeout=3s 1.1.1.1 123

Программа должна подключаться к указанному хосту (ip или доменное имя) и порту по протоколу TCP.
После подключения STDIN программы должен записываться в сокет, а данные полученные и сокета должны выводиться в STDOUT
Опционально в программу можно передать таймаут на подключение к серверу (через аргумент --timeout, по умолчанию 10s).

При нажатии Ctrl+D программа должна закрывать сокет и завершаться. Если сокет закрывается со стороны сервера, программа должна также завершаться.
При подключении к несуществующему сервер, программа должна завершаться через timeout.
*/

var timeout = flag.Int("timeout", 10, "таймаут подключения к серверу")

// write - функция записи в сокет (считываем stdin и записываем в сокет)
func write(conn net.Conn) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(conn.LocalAddr().String() + "> ")
		cmd, err := reader.ReadString('\n')
		if err != nil {
			conn.Close()
			log.Fatal(err)
		}
		_, err = conn.Write([]byte(cmd))
		if err != nil {
			conn.Close()
			log.Fatal(err)
		}
	}
}

// read - функция чтения из сокета
func read(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		msg := make([]byte, 1024)
		n, err := reader.Read(msg)
		if err != nil {
			conn.Close()
			log.Fatal(err)
		}
		fmt.Println("\n" + conn.RemoteAddr().String() + ": " + string(msg[:n]))
		fmt.Print(conn.LocalAddr().String() + "> ")
	}
}

func main() {
	// парсим таймаут и хост с портом
	flag.Parse()
	if len(flag.Args()) < 2 {
		log.Fatal("укажите хост и порт")
	}

	host := flag.Arg(0)
	port := flag.Arg(1)
	t := time.Duration(*timeout) * time.Second

	// устанавливаем соединение
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), t)
	if err != nil {
		log.Fatal(err)
	}

	// канал для сигнала ctrl+D
	signalChannel := make(chan os.Signal)
	signal.Notify(signalChannel, syscall.SIGQUIT)

	// запускаем функции записи и чтения
	go write(conn)
	go read(conn)

	// обработка ctrl+D
	select {
	case <-signalChannel:
		conn.Close()
	}
}
