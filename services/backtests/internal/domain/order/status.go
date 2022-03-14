package order

type Status string

const (
	StatusIsPending         Status = "pending"
	StatusIsPartiallyFilled Status = "partially-filled"
	StatusIsFilled          Status = "filled"
)

var (
	Statuses = []Status{
		StatusIsPending,
		StatusIsPartiallyFilled,
		StatusIsFilled,
	}
)

func (s Status) Validate() error {
	for _, vs := range Statuses {
		if s == vs {
			return nil
		}
	}

	return ErrInvalidSide
}

func (s Status) String() string {
	return string(s)
}
