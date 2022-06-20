package status

import "encoding/json"

type Status struct {
	Finished bool `json:"finished"`
}

func (s Status) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

func (s *Status) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, s)
}
