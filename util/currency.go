package util

const (
	USD = "USD"
	EUR = "EUR"
	BRL = "BRL"
)

func IsSupportedCurrencies(currency string) bool {
	switch currency {
	case USD, EUR, BRL:
		return true
	}
	return false
}
