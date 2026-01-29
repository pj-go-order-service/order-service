package order

type ProductID string

type Item struct {
	ProductID ProductID
	Name      string
	Price     Money
	Quantity  int
}
