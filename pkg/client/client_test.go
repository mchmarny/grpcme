package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	t.Run("new client sans address", func(t *testing.T) {
		if _, err := NewClient(""); err == nil {
			t.Fatalf("expected error without address")
		}
	})

	t.Run("close client", func(t *testing.T) {
		c := &SimpleClient{}
		if err := c.Close(); err != nil {
			t.Fatalf("expected no error on close")
		}
		assert.Empty(t, c.GetTarget())
	})
}
