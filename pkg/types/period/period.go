package period

import (
	"errors"
	"time"

	"golang.org/x/xerrors"
)

var (
	ErrInvalidPeriod = errors.New("invalid period")
)

type Period int64

const (
	// M1 represents 1 minute in epoch time
	M1 Period = 60
	// M3 represents 3 minutes in epoch time
	M3 Period = 3 * 60
	// M5 represents 5 minutes in epoch time
	M5 Period = 5 * 60
	// M15 represents 15 minutes in epoch time
	M15 Period = 15 * 60
	// M30 represents 30 minutes in epoch time
	M30 Period = 30 * 60
	// H1 represents 1 hour in epoch time
	H1 Period = 1 * 60 * 60
	// H2 represents 2 hours in epoch time
	H2 Period = 2 * 60 * 60
	// H4 represents 4 hours in epoch time
	H4 Period = 4 * 60 * 60
	// H6 represents 6 hours in epoch time
	H6 Period = 6 * 60 * 60
	// H8 represents 8 hours in epoch time
	H8 Period = 8 * 60 * 60
	// H12 represents 12 hours in epoch time
	H12 Period = 12 * 60 * 60
	// D1 represents 1 day in epoch time
	D1 Period = 1 * 24 * 60 * 60
	// D3 represents 3 days in epoch time
	D3 Period = 3 * 24 * 60 * 60
	// W1 represents 1 week in epoch time
	W1 Period = 1 * 7 * 24 * 60 * 60
)

func (p Period) Duration() time.Duration {
	return time.Duration(int64(p) * int64(time.Second))
}

func FromSymbol(symbol string) (Period, error) {
	switch symbol {
	case "M1":
		return M1, nil
	case "M3":
		return M3, nil
	case "M5":
		return M5, nil
	case "M15":
		return M15, nil
	case "M30":
		return M30, nil
	case "H1":
		return H1, nil
	case "H2":
		return H2, nil
	case "H4":
		return H4, nil
	case "H6":
		return H6, nil
	case "H8":
		return H8, nil
	case "H12":
		return H12, nil
	case "D1":
		return D1, nil
	case "D3":
		return D3, nil
	case "W1":
		return W1, nil
	default:
		return Period(0), xerrors.Errorf("parsing period from name (%s): %w", ErrInvalidPeriod)
	}
}

func (p Period) Name() (string, error) {
	name := p.String()
	if name == ErrInvalidPeriod.Error() {
		return "", ErrInvalidPeriod
	}
	return name, nil
}

func (p Period) String() string {
	switch p {
	case M1:
		return "M1"
	case M3:
		return "M3"
	case M5:
		return "M5"
	case M15:
		return "M15"
	case M30:
		return "M30"
	case H1:
		return "H1"
	case H2:
		return "H2"
	case H4:
		return "H4"
	case H6:
		return "H6"
	case H8:
		return "H8"
	case H12:
		return "H12"
	case D1:
		return "D1"
	case D3:
		return "D3"
	case W1:
		return "W1"
	default:
		return ErrInvalidPeriod.Error()
	}
}

func Periods() []Period {
	return []Period{
		M1, M3, M5, M15, M30,
		H1, H2, H4, H6, H8, H12,
		D1, D3,
		W1,
	}
}

func RoundTime(t time.Time, p Period) time.Time {
	diff := t.Unix() % int64(p)
	return time.Unix(t.Unix()-diff, 0)
}

func (p Period) IsAligned(t time.Time) bool {
	return (t.Unix() % int64(p)) == 0
}

func FromSeconds(i int64) (Period, error) {
	for _, p := range Periods() {
		if int64(p) == i {
			return Period(i), nil
		}
	}

	return Period(0), xerrors.Errorf("parsing period from seconds (%s): %w", ErrInvalidPeriod)
}

func (p Period) CountBetweenTimes(start, end time.Time) int64 {
	roundedStart := RoundTime(start, p)
	roundedEnd := RoundTime(end, p)
	return (roundedEnd.Unix() - roundedStart.Unix()) / int64(p)
}

func UniqueArray(per1, per2 []Period) []Period {
	tmp := make([]Period, len(per1))
	copy(tmp, per1)

	for _, p2 := range per2 {
		present := false
		for _, p1 := range per1 {
			if p1 == p2 {
				present = true
				break
			}
		}

		if !present {
			tmp = append(tmp, p2)
		}
	}

	return tmp
}

func (p Period) Opt() *Period {
	return &p
}
