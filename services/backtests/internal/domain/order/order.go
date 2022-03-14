package order

type Order struct {
	ID           uint
	Type         Type
	ExchangeName string
	PairSymbol   string
	Side         Side
	Quantity     float64
	Status       Status
}
