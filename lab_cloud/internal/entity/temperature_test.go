package entity

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCelsiusToFahrenheit(t *testing.T) {
	expected := 32.0
	input := 0.0
	actual := CelsiusToFahrenheit(input)
	assert.Equal(t, expected, actual)

	expected = -148.0
	input = -100.0
	actual = CelsiusToFahrenheit(input)
	assert.Equal(t, expected, actual)

	expected = 212.0
	input = 100.0
	actual = CelsiusToFahrenheit(input)
	assert.Equal(t, expected, actual)

}

func TestCelsiusToKelvin(t *testing.T) {
	expected := 273.0
	input := 0.0
	actual := CelsiusToKelvin(input)
	assert.Equal(t, expected, actual)

	expected = 173.0
	input = -100.0
	actual = CelsiusToKelvin(input)
	assert.Equal(t, expected, actual)

	expected = 373.0
	input = 100.0
	actual = CelsiusToKelvin(input)
	assert.Equal(t, expected, actual)
}
