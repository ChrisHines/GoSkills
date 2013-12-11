package numerics

import (
	. "github.com/smartystreets/goconvey/convey"
	"math"
	"testing"
)

func TestInvalidRange(t *testing.T) {
	Convey("Constructing a range with inverted bounds should panic", t, func() {
		So(func() {
			NewRange(10, 9)
		}, ShouldPanic)
	})
}

func TestRange(t *testing.T) {
	Convey("Given a range from 1 to 5", t, func() {
		r := NewRange(1, 5)
		Convey("0 should not be in the range", func() {
			So(r.In(0), ShouldBeFalse)
		})
		Convey("1 should be in the range", func() {
			So(r.In(1), ShouldBeTrue)
		})
		Convey("3 should be in the range", func() {
			So(r.In(3), ShouldBeTrue)
		})
		Convey("5 should be in the range", func() {
			So(r.In(5), ShouldBeTrue)
		})
		Convey("6 should not be in the range", func() {
			So(r.In(6), ShouldBeFalse)
		})
	})
}

func TestAtLeast(t *testing.T) {
	Convey("Given a range AtLeast 10", t, func() {
		r := AtLeast(10)
		Convey("0 should not be in the range", func() {
			So(r.In(0), ShouldBeFalse)
		})
		Convey("9 should not be in the range", func() {
			So(r.In(9), ShouldBeFalse)
		})
		Convey("10 should be in the range", func() {
			So(r.In(10), ShouldBeTrue)
		})
		Convey("15 should be in the range", func() {
			So(r.In(15), ShouldBeTrue)
		})
		Convey("math.MaxInt32 should be in the range", func() {
			So(r.In(math.MaxInt32), ShouldBeTrue)
		})
	})
}

func TestExactly(t *testing.T) {
	Convey("Given a range Exactly 10", t, func() {
		r := Exactly(10)
		Convey("9 should not be in the range", func() {
			So(r.In(9), ShouldBeFalse)
		})
		Convey("10 should be in the range", func() {
			So(r.In(10), ShouldBeTrue)
		})
		Convey("11 should not be in the range", func() {
			So(r.In(11), ShouldBeFalse)
		})
	})
}
