package helper

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRandomString(t *testing.T) {
	data := []struct {
		n int
	}{
		{10}, {20}, {30}, {40},
	}
	for _, s := range data {
		assert.Len(t, RandomString(s.n), s.n)
	}
}

func BenchmarkRandomString(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RandomString(100)
	}
}
