package main

import (
	"testing"
)

func BenchmarkFlipA(b *testing.B) {
	pics := readPictures()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		flipAllA(pics)
	}
}
func BenchmarkFlipB(b *testing.B) {
	pics := readPictures()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		flipAllB(pics)
	}
}
func BenchmarkFlipC(b *testing.B) {
	pics := readPictures()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		flipAllC(pics)
	}
}

func BenchmarkFlipD(b *testing.B) {
	pics := readPictures()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		flipAllD(pics)
	}
}
func BenchmarkFlipE(b *testing.B) {
	pics := readPictures()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		flipAllE(pics)
	}

}
func BenchmarkFlipF(b *testing.B) {
	pics := readPictures()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		flipAllF(pics)
	}
}

// func BenchmarkFlipG(b *testing.B) {
// 	pics := readPictures()
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		flipAllG(pics)
// 	}
// }
