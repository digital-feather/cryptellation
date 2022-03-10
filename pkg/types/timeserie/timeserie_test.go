package timeserie

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

func TestTimeSerieSuite(t *testing.T) {
	suite.Run(t, new(TimeSerieSuite))
}

type TimeSerieSuite struct {
	suite.Suite
}

func (suite *TimeSerieSuite) TestNew() {
	ts := New()
	suite.Require().Len(ts.data, 0)
	suite.Require().Len(ts.orderedKeys, 0)
}

func (suite *TimeSerieSuite) TestSet() {
	ts := New()

	ts2 := ts.Set(time.Unix(1, 0), "new")
	suite.Require().Equal(ts, ts2)
	suite.Require().Len(ts.data, 1)
	suite.Require().Len(ts.orderedKeys, 1)

	ts.Set(time.Unix(0, 0), "first")
	suite.Require().Len(ts.data, 2)
	suite.Require().Len(ts.orderedKeys, 2)
	suite.Require().Equal(time.Unix(0, 0), ts.orderedKeys[0].ToTime())
	suite.Require().Equal(time.Unix(1, 0), ts.orderedKeys[1].ToTime())

	ts.Set(time.Unix(3, 0), "last")
	suite.Require().Len(ts.data, 3)
	suite.Require().Len(ts.orderedKeys, 3)
	suite.Require().Equal(time.Unix(0, 0), ts.orderedKeys[0].ToTime())
	suite.Require().Equal(time.Unix(1, 0), ts.orderedKeys[1].ToTime())
	suite.Require().Equal(time.Unix(3, 0), ts.orderedKeys[2].ToTime())

	ts.Set(time.Unix(2, 0), "middle")
	suite.Require().Len(ts.data, 4)
	suite.Require().Len(ts.orderedKeys, 4)
	suite.Require().Equal(time.Unix(0, 0), ts.orderedKeys[0].ToTime())
	suite.Require().Equal(time.Unix(1, 0), ts.orderedKeys[1].ToTime())
	suite.Require().Equal(time.Unix(2, 0), ts.orderedKeys[2].ToTime())
	suite.Require().Equal(time.Unix(3, 0), ts.orderedKeys[3].ToTime())

	ts.Set(time.Unix(1, 0), "existing")
	suite.Require().Len(ts.data, 4)
	suite.Require().Len(ts.orderedKeys, 4)
	suite.Require().Equal(time.Unix(0, 0), ts.orderedKeys[0].ToTime())
	suite.Require().Equal(time.Unix(1, 0), ts.orderedKeys[1].ToTime())
	suite.Require().Equal(time.Unix(2, 0), ts.orderedKeys[2].ToTime())
	suite.Require().Equal(time.Unix(3, 0), ts.orderedKeys[3].ToTime())
}

func (suite *TimeSerieSuite) TestGet() {
	ts := New()
	for i := 0; i < 100; i++ {
		ts.Set(time.Unix(int64(i), 0), fmt.Sprintf("element-%d", i))
	}

	data, exists := ts.Get(time.Unix(25, 0))
	suite.Require().True(exists)
	suite.Require().Equal("element-25", data)

	data, exists = ts.Get(time.Unix(int64(200), 0))
	suite.Require().False(exists)
	suite.Require().Nil(data)
}

func (suite *TimeSerieSuite) TestLen() {
	ts := New()
	for i := 0; i < 100; i++ {
		suite.Require().Equal(i, ts.Len())
		ts.Set(time.Unix(int64(i), 0), fmt.Sprintf("element-%d", i))
	}
}

func (suite *TimeSerieSuite) TestMerge() {
	ts, ts2 := New(), New()
	for i := 0; i < 10; i++ {
		ts.Set(time.Unix(int64(i), 0), fmt.Sprintf("element-%d", i))
		ts2.Set(time.Unix(int64(i)+10, 0), fmt.Sprintf("element-%d", i+10))
	}
	suite.Require().NoError(ts.Merge(*ts2, nil))

	suite.Require().Equal(20, ts.Len())
	for i := 0; i < 20; i++ {
		suite.Require().Equal(time.Unix(int64(i), 0), ts.orderedKeys[i].ToTime())
		suite.Require().Equal(fmt.Sprintf("element-%d", i), ts.data[newKey(time.Unix(int64(i), 0))])
	}
}

func (suite *TimeSerieSuite) TestMergeWithoutCollision() {
	ts, ts2 := New(), New()
	for i := 0; i < 10; i++ {
		ts.Set(time.Unix(int64(i), 0), fmt.Sprintf("element-%d", i))
		ts2.Set(time.Unix(int64(i+10), 0), fmt.Sprintf("element-%d", i+10))
	}
	suite.Require().NoError(ts.Merge(*ts2, &MergeOptions{
		ErrorOnCollision: true,
	}))

	suite.Require().Equal(20, ts.Len())
	for i := 0; i < 20; i++ {
		suite.Require().Equal(time.Unix(int64(i), 0), ts.orderedKeys[i].ToTime())
		suite.Require().Equal(fmt.Sprintf("element-%d", i), ts.data[newKey(time.Unix(int64(i), 0))])
	}
}

func (suite *TimeSerieSuite) TestMergeWithCollision() {
	ts, ts2 := New(), New()
	for i := 0; i < 10; i++ {
		ts.Set(time.Unix(int64(i), 0), fmt.Sprintf("element-%d", i))
		ts2.Set(time.Unix(int64(i)+5, 0), fmt.Sprintf("element-%d", i+10))
	}
	suite.Require().Error(ts.Merge(*ts2, &MergeOptions{
		ErrorOnCollision: true,
	}))

	suite.Require().Equal(10, ts.Len())
	for i := 0; i < 10; i++ {
		suite.Require().Equal(time.Unix(int64(i), 0), ts.orderedKeys[i].ToTime())
		suite.Require().Equal(fmt.Sprintf("element-%d", i), ts.data[newKey(time.Unix(int64(i), 0))])
	}
}

func (suite *TimeSerieSuite) TestDelete() {
	ts := New()
	for i := 0; i < 100; i++ {
		ts.Set(time.Unix(int64(i), 0), fmt.Sprintf("element-%d", i))
	}

	ts.Delete(time.Unix(99, 0))
	suite.Require().Equal(99, ts.Len())
	for i := 0; i < 99; i++ {
		suite.Require().Equal(time.Unix(int64(i), 0), ts.orderedKeys[i].ToTime())
		suite.Require().Equal(fmt.Sprintf("element-%d", i), ts.data[newKey(time.Unix(int64(i), 0))])
	}
	_, exists := ts.data[newKey(time.Unix(99, 0))]
	suite.Require().False(exists)

	ts.Delete(time.Unix(0, 0))
	suite.Require().Equal(98, ts.Len())
	for i := 0; i < 98; i++ {
		suite.Require().Equal(time.Unix(int64(i)+1, 0), ts.orderedKeys[i].ToTime())
		suite.Require().Equal(fmt.Sprintf("element-%d", i+1), ts.data[newKey(time.Unix(int64(i)+1, 0))])
	}
	_, exists = ts.data[newKey(time.Unix(0, 0))]
	suite.Require().False(exists)

	ts.Delete(time.Unix(50, 0))
	suite.Require().Equal(97, ts.Len())
	for i := 0; i < 97; i++ {
		if i+1 < 50 {
			suite.Require().Equal(time.Unix(int64(i)+1, 0), ts.orderedKeys[i].ToTime())
			suite.Require().Equal(fmt.Sprintf("element-%d", i+1), ts.data[newKey(time.Unix(int64(i)+1, 0))])
		} else {
			suite.Require().Equal(time.Unix(int64(i)+2, 0), ts.orderedKeys[i].ToTime())
			suite.Require().Equal(fmt.Sprintf("element-%d", i+2), ts.data[newKey(time.Unix(int64(i)+2, 0))])
		}
	}
	_, exists = ts.data[newKey(time.Unix(50, 0))]
	suite.Require().False(exists)
}

type loopTestObject struct {
	Time time.Time
	Obj  interface{}
}

func (suite *TimeSerieSuite) TestLoop() {
	ts := New()

	ts.Set(time.Unix(0, 0), "zero")
	ts.Set(time.Unix(1, 0), "un")
	ts.Set(time.Unix(2, 0), "deux")

	tsList := []loopTestObject{}
	suite.Require().NoError(ts.Loop(func(ts time.Time, obj interface{}) (bool, error) {
		tsList = append(tsList, loopTestObject{
			Time: ts,
			Obj:  obj,
		})
		return false, nil
	}))

	suite.Require().Len(tsList, 3)
	suite.Require().Equal(tsList[0].Time, time.Unix(0, 0))
	suite.Require().Equal(tsList[0].Obj, "zero")
	suite.Require().Equal(tsList[1].Time, time.Unix(1, 0))
	suite.Require().Equal(tsList[1].Obj, "un")
	suite.Require().Equal(tsList[2].Time, time.Unix(2, 0))
	suite.Require().Equal(tsList[2].Obj, "deux")
}

func (suite *TimeSerieSuite) TestLoopBreak() {
	ts := New()

	ts.Set(time.Unix(0, 0), "zero")
	ts.Set(time.Unix(1, 0), "un")
	ts.Set(time.Unix(2, 0), "deux")

	tsList := []loopTestObject{}
	suite.Require().NoError(ts.Loop(func(ts time.Time, obj interface{}) (bool, error) {
		tsList = append(tsList, loopTestObject{
			Time: ts,
			Obj:  obj,
		})

		if len(tsList) == 1 {
			return true, nil
		}

		return false, nil
	}))

	suite.Require().Len(tsList, 1)
	suite.Require().Equal(tsList[0].Time, time.Unix(0, 0))
	suite.Require().Equal(tsList[0].Obj, "zero")
}

func (suite *TimeSerieSuite) TestLoopError() {
	ts := New()

	ts.Set(time.Unix(0, 0), "zero")
	ts.Set(time.Unix(1, 0), "un")
	ts.Set(time.Unix(2, 0), "deux")

	tsList := []loopTestObject{}
	suite.Require().Error(ts.Loop(func(ts time.Time, obj interface{}) (bool, error) {
		tsList = append(tsList, loopTestObject{
			Time: ts,
			Obj:  obj,
		})

		if len(tsList) == 1 {
			return true, errors.New("test-error")
		}

		return false, nil
	}))

	suite.Require().Len(tsList, 1)
	suite.Require().Equal(tsList[0].Time, time.Unix(0, 0))
	suite.Require().Equal(tsList[0].Obj, "zero")
}

func (suite *TimeSerieSuite) TestFirst() {
	ts := New()

	_, _, ok := ts.First()
	suite.Require().False(ok)

	ts.Set(time.Unix(1, 0), "new")
	t, l, ok := ts.First()
	suite.Require().True(ok)
	suite.Require().Equal(time.Unix(1, 0), t)
	suite.Require().Equal("new", l)

	ts.Set(time.Unix(0, 0), "first")
	t, l, ok = ts.First()
	suite.Require().True(ok)
	suite.Require().Equal(time.Unix(0, 0), t)
	suite.Require().Equal("first", l)

	ts.Set(time.Unix(3, 0), "last")
	t, l, ok = ts.First()
	suite.Require().True(ok)
	suite.Require().Equal(time.Unix(0, 0), t)
	suite.Require().Equal("first", l)
}

func (suite *TimeSerieSuite) TestLast() {
	ts := New()

	_, _, ok := ts.Last()
	suite.Require().False(ok)

	ts.Set(time.Unix(1, 0), "new")
	t, l, ok := ts.Last()
	suite.Require().True(ok)
	suite.Require().Equal(time.Unix(1, 0), t)
	suite.Require().Equal("new", l)

	ts.Set(time.Unix(0, 0), "first")
	t, l, ok = ts.Last()
	suite.Require().True(ok)
	suite.Require().Equal(time.Unix(1, 0), t)
	suite.Require().Equal("new", l)

	ts.Set(time.Unix(3, 0), "last")
	t, l, ok = ts.Last()
	suite.Require().True(ok)
	suite.Require().Equal(time.Unix(3, 0), t)
	suite.Require().Equal("last", l)

	ts.Set(time.Unix(2, 0), "middle")
	t, l, ok = ts.Last()
	suite.Require().True(ok)
	suite.Require().Equal(time.Unix(3, 0), t)
	suite.Require().Equal("last", l)

	ts.Set(time.Unix(1, 0), "existing")
	t, l, ok = ts.Last()
	suite.Require().True(ok)
	suite.Require().Equal(time.Unix(3, 0), t)
	suite.Require().Equal("last", l)
}

func (suite *TimeSerieSuite) TestExtract() {
	ts := New()
	for i := int64(0); i < 4; i++ {
		ts.Set(time.Unix(60*i, 0), i)
	}

	nl := ts.Extract(time.Unix(60, 0), time.Unix(120, 0))
	suite.Require().Equal(2, nl.Len())

	obj, exists := nl.Get(time.Unix(60, 0))
	suite.Require().True(exists)
	suite.Require().Equal(int64(1), obj)

	obj, exists = nl.Get(time.Unix(120, 0))
	suite.Require().True(exists)
	suite.Require().Equal(int64(2), obj)
}

func (suite *TimeSerieSuite) TestFirstN() {
	ts := New()
	for i := int64(0); i < 4; i++ {
		ts.Set(time.Unix(60*i, 0), i)
	}

	nl := ts.FirstN(2)
	suite.Require().Equal(2, nl.Len())

	obj, exists := nl.Get(time.Unix(0, 0))
	suite.Require().True(exists)
	suite.Require().Equal(int64(0), obj)

	obj, exists = nl.Get(time.Unix(60, 0))
	suite.Require().True(exists)
	suite.Require().Equal(int64(1), obj)
}
