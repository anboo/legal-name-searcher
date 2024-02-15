package main

import (
	"os"
	"testing"
)

func BenchmarkCountLines(b *testing.B) {
	filename := "data-20220314-structure-20220314.csv"

	b.Run("Scanner", func(b *testing.B) {
		f, _ := os.Open("filename")
		for i := 0; i < b.N; i++ {
			_, err := countLinesWithScanner(f)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("WC", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := countLinesWithWC(filename)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("syscalls", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := countLinesWithSyscalls(filename)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
