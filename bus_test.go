package event_bus

import (
	"github.com/google/uuid"
	"testing"
)

func TestAll(t *testing.T) {
	t.Run("GenerateUUID", func(t *testing.T) {
		t.Log(uuid.New().String())
	})
}
