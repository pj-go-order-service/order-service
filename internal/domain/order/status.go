package order

type Status string

const (
	StatusCreated  Status = "created"
	StatusPaid     Status = "paid"
	StatusCanceled Status = "canceled"
)
