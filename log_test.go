package chata_test

import (
	"context"
	"log/slog"
	"testing"

	"github.com/rchamarthy/chata"
	"github.com/stretchr/testify/assert"
)

func TestLog(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	assert := assert.New(t)

	log := chata.Log(ctx)
	assert.NotNil(log)
	assert.Equal(log, chata.NilLogger())

	ctx = context.WithValue(ctx, chata.LogKey, slog.Default())
	log = chata.Log(ctx)
	assert.NotNil(log)
	assert.NotEqual(log, chata.NilLogger())
	assert.Equal(log, slog.Default())

}

func TestNilLogger(t *testing.T) {
	t.Parallel()

	nilLogger := chata.NilLogger()
	assert.NotNil(t, nilLogger)
}
