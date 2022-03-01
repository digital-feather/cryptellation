package period

import (
	"testing"
	"time"

	"golang.org/x/xerrors"
)

func TestPeriodDuration(t *testing.T) {
	if M1.Duration() != time.Minute {
		t.Error("Period and duration mismatched:", M1, time.Minute)
	}
}

func TestRoundTime(t *testing.T) {
	toCorrect, err := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
	if err != nil {
		t.Fatal(err)
	}

	goal, err := time.Parse(time.RFC3339, "2006-01-02T15:04:00Z")
	if err != nil {
		t.Fatal(err)
	}

	corrected := M1.RoundTime(toCorrect)
	if !goal.Equal(corrected) {
		t.Error("These two should be equal:", goal, corrected)
	}
}

func TestPeriodSymbols(t *testing.T) {
	symbols := Symbols()

	if len(symbols) != 14 {
		t.Error("There should be 14 period symbols but there is", len(symbols))
	}

	if !inArray(symbols, M1) {
		t.Error("There is no M1")
	} else if !inArray(symbols, M3) {
		t.Error("There is no M3")
	} else if !inArray(symbols, M5) {
		t.Error("There is no M5")
	} else if !inArray(symbols, M15) {
		t.Error("There is no M15")
	} else if !inArray(symbols, M30) {
		t.Error("There is no M30")
	} else if !inArray(symbols, H1) {
		t.Error("There is no H1")
	} else if !inArray(symbols, H2) {
		t.Error("There is no H2")
	} else if !inArray(symbols, H4) {
		t.Error("There is no H4")
	} else if !inArray(symbols, H6) {
		t.Error("There is no H6")
	} else if !inArray(symbols, H8) {
		t.Error("There is no H8")
	} else if !inArray(symbols, H12) {
		t.Error("There is no H12")
	} else if !inArray(symbols, D1) {
		t.Error("There is no D1")
	} else if !inArray(symbols, D3) {
		t.Error("There is no D3")
	} else if !inArray(symbols, W1) {
		t.Error("There is no W1")
	}
}

func TestSymbolsString(t *testing.T) {
	for _, s := range Symbols() {
		if s.String() == ErrInvalidPeriod.Error() {
			t.Error("There is no string for", s)
		}
	}
}

func TestValidateSymbol(t *testing.T) {
	for _, s := range Symbols() {
		if s.Validate() != nil {
			t.Errorf("There is an error for %q", s)
		}
	}

	wrong := Symbol("unknown")
	if !xerrors.Is(wrong.Validate(), ErrInvalidPeriod) {
		t.Errorf("Wrong symbol should be %q but is %q", ErrInvalidPeriod, wrong.Validate())
	}
}

func TestIsAligned(t *testing.T) {
	if !M1.IsAligned(time.Unix(60, 0)) {
		t.Error("Time 60 should be aligned on M1")
	}

	if M1.IsAligned(time.Unix(45, 0)) {
		t.Error("Time 45 should not be aligned on M1")
	}
}

func inArray(array []Symbol, element Symbol) bool {
	for _, k := range array {
		if k == element {
			return true
		}
	}
	return false
}

func TestFromSeconds(t *testing.T) {
	if _, err := FromSeconds(60); err != nil {
		t.Error("There should be no error for 60s")
	}

	if _, err := FromSeconds(59); err == nil {
		t.Error("There should be an error for 59s")
	}
}

func TestCountBetweenTimes(t *testing.T) {
	symbs := []Symbol{M1, D1}

	for _, s := range symbs {
		now := time.Now()
		count := s.CountBetweenTimes(now, now)
		if count != 0 {
			t.Error("Count between times should be 0 between same dates")
		}

		count = s.CountBetweenTimes(now.Add(-s.Duration()), now)
		if count != 1 {
			t.Errorf("Count between times should be 1 between dates that are just one period apart but is %d", count)
		}

		count = s.CountBetweenTimes(now.Add(-s.Duration()*10), now)
		if count != 10 {
			t.Errorf("Count between times should be 10 between dates that are just ten period apart but is %d", count)
		}
	}
}

func TestUniqueArray(t *testing.T) {
	s1 := []Symbol{M1, M15}
	s2 := []Symbol{M1, M3}
	s3 := []Symbol{M1, M3, M15}

	m := UniqueArray(s2, s1)
	if len(m) != 3 || m[0] != s3[0] || m[1] != s3[1] || m[2] != s3[2] {
		t.Error(s3, m)
	}
}
