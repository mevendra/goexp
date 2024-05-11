package entity

func CelsiusToKelvin(temp float64) float64 {
	return temp + 273
}

func CelsiusToFahrenheit(temp float64) float64 {
	return temp*1.8 + 32
}
