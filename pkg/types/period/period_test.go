package period

import (
	"testing"
	"time"
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

	corrected := RoundTime(toCorrect, M1)

	if !goal.Equal(corrected) {
		t.Error("These two should be equal:", goal, corrected)
	}
}

func TestPeriods(t *testing.T) {
	periods := Periods()

	if len(periods) != 14 {
		t.Error("There should be 14 periods but there is", len(periods))
	}

	if !inArray(periods, M1) {
		t.Error("There is no M1")
	} else if !inArray(periods, M3) {
		t.Error("There is no M3")
	} else if !inArray(periods, M5) {
		t.Error("There is no M5")
	} else if !inArray(periods, M15) {
		t.Error("There is no M15")
	} else if !inArray(periods, M30) {
		t.Error("There is no M30")
	} else if !inArray(periods, H1) {
		t.Error("There is no H1")
	} else if !inArray(periods, H2) {
		t.Error("There is no H2")
	} else if !inArray(periods, H4) {
		t.Error("There is no H4")
	} else if !inArray(periods, H6) {
		t.Error("There is no H6")
	} else if !inArray(periods, H8) {
		t.Error("There is no H8")
	} else if !inArray(periods, H12) {
		t.Error("There is no H12")
	} else if !inArray(periods, D1) {
		t.Error("There is no D1")
	} else if !inArray(periods, D3) {
		t.Error("There is no D3")
	} else if !inArray(periods, W1) {
		t.Error("There is no W1")
	}
}

func TestString(t *testing.T) {
	for _, p := range Periods() {
		if p.String() == ErrInvalidPeriod.Error() {
			t.Error("There is no string for", p)
		}
	}

	wrong := Period(123)
	if wrong.String() != ErrInvalidPeriod.Error() {
		t.Error("Wrong period should be", ErrInvalidPeriod.Error(), " but is", wrong.String())
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

func inArray(array []Period, element Period) bool {
	for _, k := range array {
		if k == element {
			return true
		}
	}
	return false
}

func TestFromInt(t *testing.T) {
	for _, c := range Periods() {
		if _, err := FromSeconds(int64(c)); err != nil {
			t.Error("There should be no error for", c)
		}
	}
}

func TestFromIntError(t *testing.T) {
	if _, err := FromSeconds(0); err == nil {
		t.Error("There should be an error")
	}
}

func TestFromName(t *testing.T) {
	// TODO
}

func TestName(t *testing.T) {
	// TODO
}

func TestCountBetweenTimes(t *testing.T) {
	pers := []Period{M1, D1}

	for _, p := range pers {
		now := time.Now()
		count := p.CountBetweenTimes(now, now)
		if count != 0 {
			t.Error("Count between times should be 0 between same dates")
		}

		count = p.CountBetweenTimes(now.Add(-p.Duration()), now)
		if count != 1 {
			t.Error("Count between times should be 1 between dates that are just one period apart")
		}

		count = p.CountBetweenTimes(now.Add(-p.Duration()*10), now)
		if count != 10 {
			t.Error("Count between times should be 10 between dates that are just ten period apart")
		}
	}
}

func TestUniqueArray(t *testing.T) {
	p1 := []Period{M1, M15}
	p2 := []Period{M1, M3}
	p3 := []Period{M1, M3, M15}

	m := UniqueArray(p2, p1)
	if len(m) != 3 || m[0] != p3[0] || m[1] != p3[1] || m[2] != p3[2] {
		t.Error(p3, m)
	}
}
