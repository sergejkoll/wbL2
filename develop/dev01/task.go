package main

import (
	"fmt"
	"log"
	"time"

	"github.com/beevik/ntp"
)

/*
=== Базовая задача ===

Создать программу печатающую точное время с использованием NTP библиотеки.Инициализировать как go module.
Использовать библиотеку https://github.com/beevik/ntp.
Написать программу печатающую текущее время / точное время с использованием этой библиотеки.

Программа должна быть оформлена с использованием как go module.
Программа должна корректно обрабатывать ошибки библиотеки: распечатывать их в STDERR и возвращать ненулевой код выхода в OS.
Программа должна проходить проверки go vet и golint.
*/

func main() {
	// Получние текущего времени от NTP-сервера
	now, err := ntp.Time("0.beevik-ntp.pool.ntp.org")
	if err != nil {
		// Fatal - вызывает os.Exit(1) а также log записывает сообщения в stderr
		log.Fatal(err)
	}
	fmt.Println(now)

	// Запрос с получением метаданных
	resp, err := ntp.Query("0.beevik-ntp.pool.ntp.org")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp.Time.Local())
	fmt.Println(resp.ClockOffset)
	// Вывод точного времени со смещением
	fmt.Println(time.Now().Add(resp.ClockOffset))
}
