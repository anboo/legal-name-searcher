package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func countLinesWithScanner(file *os.File) (int, error) {
	scanner := bufio.NewScanner(file)
	count := 0

	for scanner.Scan() {
		count++
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return count, nil
}

func countLinesWithWC(filename string) (int, error) {
	cmd := exec.Command("wc", "-l", filename)
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	fmt.Println(string(output))
	return 0, nil
}

func countLinesWithSyscalls(filename string) (int, error) {
	// Открываем файл
	fd, err := syscall.Open(filename, syscall.O_RDONLY, 0)
	if err != nil {
		return 0, err
	}
	defer syscall.Close(fd)

	// Буфер для чтения данных
	buffer := make([]byte, 1024)
	count := 0

	for {
		// Читаем блок данных
		n, err := syscall.Read(fd, buffer)
		if err != nil || n == 0 {
			break
		}

		// Подсчитываем символы новой строки
		for i := 0; i < n; i++ {
			if buffer[i] == '\n' {
				count++
			}
		}
	}

	return count, nil
}
