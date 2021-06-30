//[+]Суть задания — написать программу, которая считывает из файла список
//математических выражений, считает результат и записывает в другой файл.
//Пример входного файла:
//5+4=?
// 9+3=?
// Сегодня прекрасная погода
// 13+7=?
// 4-2=?
// Пример файла с выводом:
// 5+4=9
// 9+3=12
// 13+7=20
// 4-2=2

// [+]Использовать методы и структуры пакетов ioutils и regexp.
// [+]Программа должна принимать на вход 2 аргумента: имя входного файла и имя файла для вывода результатов.
// [+]Если не найден вывод, создать.
// [+]Если файл вывода существует, очистить перед записью новых результатов.
// [+]Использовать буферизированную запись результатов.

package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
)

func main() {

	args := os.Args
	if len(args) < 3 {
		log.Fatal("Wrong args. Usage maths.exe <source file> <target file>")
	}

	sourceFile, err := os.Open(args[1])
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(args[2], []byte{}, 0777)
	if err != nil {
		log.Fatal(err)
	}

	targetFile, err := os.OpenFile(args[2], os.O_APPEND, 0777)
	if err != nil {
		log.Fatal(err)
	}
	defer targetFile.Close()

	reader := bufio.NewReaderSize(sourceFile, 2) // <-- уменьшенный буфер
	writer := bufio.NewWriter(targetFile)
	exp := regexp.MustCompile(`^([0-9]+)([+-/*])([0-9]+)=\?$`)

	for {
		line, err := liner(reader) // <-- обработка префикса
		if err == io.EOF {
			writer.Flush()
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		sub := exp.FindStringSubmatch(string(line))
		if len(sub) < 1 {
			continue
		}

		a, err := strconv.Atoi(sub[1])
		if err != nil {
			log.Fatal(err)
		}
		b, err := strconv.Atoi(sub[3])
		if err != nil {
			log.Fatal(err)
		}
		_, err = writer.WriteString(fmt.Sprintf("%d%s%d=%d\n",
			a, sub[2], b, mathAction(a, b, sub[2])))
		if err != nil {
			log.Fatal(err)
		}
	}
}

func liner(r *bufio.Reader) ([]byte, error) {
	line, isPrefix, err := r.ReadLine()
	if err != nil {
		return nil, err
	}
	if isPrefix {
		newline, err := liner(r)
		if err != nil {
			return line, err
		}
		line = append(line, newline...)
	}
	return line, err
}

func mathAction(a, b int, sign string) int {
	switch sign {
	case "+":
		return a + b
	case "-":
		return a - b
	case "*":
		return a * b
	case "/":
		return a / b
	}
	return 0
}
