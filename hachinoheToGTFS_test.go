package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToTime610(t *testing.T) {
	actual := toTime("610")
	expected := "10:10:00"
	assert.Equal(t, expected, actual)
}

func TestToTime301(t *testing.T) {
	actual := toTime("301")
	expected := "05:01:00"
	assert.Equal(t, expected, actual)
}
func TestMaeZero(t *testing.T) {
	actual := maeZero("3", 7)
	expected := "0000003"
	assert.Equal(t, expected, actual)
}
