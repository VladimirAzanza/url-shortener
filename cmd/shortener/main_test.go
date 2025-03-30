package main

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
)

func TestValidateApp(t *testing.T) {
	err := fx.ValidateApp(Module)
	require.NoError(t, err)
}
