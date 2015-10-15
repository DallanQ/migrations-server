package main

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUtils(t *testing.T) {
	Convey("Foo", t, func() {
		Convey("foo test", func() {
			So(2+2, ShouldEqual, 4)
		})
	})
}

func TestFake(t *testing.T) {
	Convey("Fake", t, func() {
		So(1+1, ShouldEqual, 2)
	})
}
