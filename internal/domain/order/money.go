package order

type Money struct {
	Amount   int64 // в минимальных единицах
	Currency string
}

func NewMoney(amount int64, currency string) Money {
	return Money{
		Amount:   amount,
		Currency: currency,
	}
}
