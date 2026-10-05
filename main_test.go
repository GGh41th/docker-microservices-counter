// main_test.go
package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

func TestHandlerIncrementsCounter(t *testing.T) {
	mr := miniredis.RunT(t)
	rdb = redis.NewClient(&redis.Options{Addr: mr.Addr()})
	t.Cleanup(func() { rdb.Close() })

	for i := 1; i <= 3; i++ {
		rec := httptest.NewRecorder()
		handler(rec, httptest.NewRequest(http.MethodGet, "/", nil))

		want := fmt.Sprintf("vue %d fois", i)
		if !strings.Contains(rec.Body.String(), want) {
			t.Fatalf("request %d: want %q in body, got %q", i, want, rec.Body.String())
		}
	}
}

func TestHandlerRedisDown(t *testing.T) {
	mr := miniredis.RunT(t)
	rdb = redis.NewClient(&redis.Options{Addr: mr.Addr()})
	t.Cleanup(func() { rdb.Close() })
	mr.Close() // simulate the DB going away

	rec := httptest.NewRecorder()
	handler(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	if !strings.Contains(rec.Body.String(), "erreur Redis") {
		t.Fatalf("expected Redis error message, got %q", rec.Body.String())
	}
}