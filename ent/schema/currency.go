package schema

type Currency string

func (_ Currency) Values() []string {
	return []string{"RUB", "USD", "EUR", "GBP"}
}

func (c Currency) Valid() bool {
	return c == "RUB" || c == "USD" || c == "EUR" || c == "GBP"
}

func (c Currency) String() string {
	return string(c)
}
