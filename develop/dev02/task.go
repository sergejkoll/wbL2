package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
)

/*
=== Задача на распаковку ===

Создать Go функцию, осуществляющую примитивную распаковку строки, содержащую повторяющиеся символы / руны, например:
	- "a4bc2d5e" => "aaaabccddddde"
	- "abcd" => "abcd"
	- "45" => "" (некорректная строка)
	- "" => ""
Дополнительное задание: поддержка escape - последовательностей
	- qwe\4\5 => qwe45 (*)
	- qwe\45 => qwe44444 (*)
	- qwe\\5 => qwe\\\\\ (*)

В случае если была передана некорректная строка функция должна возвращать ошибку. Написать unit-тесты.

Функция должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

// функция записи повторящихся символов
func writeBlock(builder *strings.Builder, number *int, char rune) {
	if *number != 0 {
		for i := 0; i < *number-1; i++ {
			builder.WriteRune(char)
		}
		*number = 0
	}
}

// функция распаковки
func unpack(str string) (result string, err error) {
	var resultB strings.Builder

	var previousRune rune
	num := 0
	escape := false
	previousIsEscape := false

	for idx, current := range str {
		// если встретили escape
		if current == '\\' {
			// если при этом escape уже был то сохраняем текущий элемент и записываем его
			if escape {
				resultB.WriteRune(current)
				previousRune = current
				escape = false
			} else {
				escape = true
			}
			continue
		}

		// если прошлый символ был escape
		if escape {
			escape = false
			resultB.WriteRune(current)
			previousRune = current
			previousIsEscape = true // для обработки ситуации \45
			continue
		}

		// преобразование руны в число если это не число записываем в строку и сохраняем
		num, err = strconv.Atoi(string(current))
		if err != nil {
			resultB.WriteRune(current)
			previousRune = current
			continue
		}

		_, err = strconv.Atoi(string(previousRune)) // проверка чтобы два символа подряд не были числами

		switch {
		case err == nil && previousIsEscape == false:
			return "", errors.New("invalid string")
		case idx == 0:
			return "", errors.New("number in zero position")
		case num == 0:
			return "", errors.New("number is zero")
		}

		writeBlock(&resultB, &num, previousRune) // запись блока

		previousRune = current
		if previousRune != '\\' {
			previousIsEscape = false
		}
	}

	result = resultB.String()
	return result, nil
}

func main() {
	res, err := unpack("a4bc2d5e")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
}
