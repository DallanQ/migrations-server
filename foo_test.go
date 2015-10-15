package main_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFoo(t *testing.T) {
	Convey("Foo", t, func() {
		Convey("foo called", func() {
			So(2+2, ShouldEqual, 4)
		})
	})
}
