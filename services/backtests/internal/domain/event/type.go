package event

type Type string

const (
	TypeIsTick   Type = "tick"
	TypeIsStatus Type = "status"
)

func (t Type) String() string {
	return string(t)
}

func (t Type) MarshalBinary() ([]byte, error) {
	return []byte(t.String()), nil
}

func (t Type) UnmarshalBinary(data []byte) error {
	t = Type(string(data))
	return nil
}
