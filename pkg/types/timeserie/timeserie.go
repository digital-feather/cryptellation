package timeserie

import (
	"errors"
	"time"
)

const (
	keyFormat = time.RFC3339Nano
)

var (
	ErrTimeStampAlreadyExists error = errors.New("timestamp-already-exists")
)

type key int64

func newKey(t time.Time) key {
	return key(t.UnixNano())
}

func (k key) ToTime() time.Time {
	return time.Unix(0, int64(k))
}

type TimeSerie struct {
	data        map[key]interface{}
	orderedKeys []key
}

func New() *TimeSerie {
	return &TimeSerie{
		data:        make(map[key]interface{}),
		orderedKeys: make([]key, 0),
	}
}

func (ts *TimeSerie) Set(t time.Time, d interface{}) *TimeSerie {
	k := newKey(t)

	ts.addKey(k)
	ts.data[k] = d

	return ts
}

func (ts *TimeSerie) addKey(k key) {
	if _, exists := ts.data[k]; exists {
		return
	}

	previous := -1
	for pos, v := range ts.orderedKeys {
		if v > k {
			break
		}
		previous = pos
	}

	if previous == -1 {
		ts.orderedKeys = append([]key{k}, ts.orderedKeys...)
	} else if previous == len(ts.orderedKeys)-1 {
		ts.orderedKeys = append(ts.orderedKeys, k)
	} else {
		index := previous + 1
		ts.orderedKeys = append(ts.orderedKeys[:index+1], ts.orderedKeys[index:]...)
		ts.orderedKeys[index] = k
	}
}

func (ts *TimeSerie) Get(t time.Time) (interface{}, bool) {
	data, exists := ts.data[newKey(t)]
	return data, exists
}

func (ts *TimeSerie) Len() int {
	return len(ts.orderedKeys)
}

type MergeOptions struct {
	ErrorOnCollision bool
}

func (ts *TimeSerie) Merge(ts2 TimeSerie, options *MergeOptions) error {
	if options != nil && options.ErrorOnCollision {
		for pos := range ts2.data {
			if _, exists := ts.Get(pos.ToTime()); exists {
				return ErrTimeStampAlreadyExists
			}
		}
	}

	for pos, v := range ts2.data {
		ts.Set(pos.ToTime(), v)
	}

	return nil
}

func (ts *TimeSerie) Delete(t ...time.Time) {
	for _, ti := range t {
		k := newKey(ti)

		if _, exists := ts.data[k]; !exists {
			return
		}

		delete(ts.data, k)
		index := -1
		for pos, v := range ts.orderedKeys {
			index = pos
			if v == k {
				break
			}
		}

		if index == 0 {
			ts.orderedKeys = ts.orderedKeys[1:]
		} else if index == len(ts.orderedKeys)-1 {
			ts.orderedKeys = ts.orderedKeys[:len(ts.orderedKeys)-1]
		} else {
			ts.orderedKeys = append(ts.orderedKeys[:index], ts.orderedKeys[index+1:]...)
		}
	}
}

func (ts *TimeSerie) Loop(callback func(time.Time, interface{}) (bool, error)) error {
	for _, t := range ts.orderedKeys {
		if shouldBreak, err := callback(t.ToTime(), ts.data[t]); err != nil {
			return err
		} else if shouldBreak {
			break
		}
	}

	return nil
}

func (ts TimeSerie) First() (time.Time, interface{}, bool) {
	if ts.Len() == 0 {
		return time.Unix(0, 0), nil, false
	}

	t := ts.orderedKeys[0]
	return t.ToTime(), ts.data[t], true
}

func (ts TimeSerie) Last() (time.Time, interface{}, bool) {
	if ts.Len() == 0 {
		return time.Unix(0, 0), nil, false
	}

	t := ts.orderedKeys[ts.Len()-1]
	return t.ToTime(), ts.data[t], true
}

func (ts TimeSerie) Extract(start, end time.Time) *TimeSerie {
	ets := New()

	ts.Loop(func(t time.Time, obj interface{}) (bool, error) {
		if t.Before(start) {
			return false, nil
		}

		if t.After(end) {
			return true, nil
		}

		ets.Set(t, obj)
		return false, nil
	})

	return ets
}

func (ts TimeSerie) FirstN(limit uint) *TimeSerie {
	ets := New()

	if limit == 0 {
		return ets
	}

	var count uint
	ts.Loop(func(t time.Time, obj interface{}) (bool, error) {
		ets.Set(t, obj)
		count++
		return count >= limit, nil
	})

	return ets
}
