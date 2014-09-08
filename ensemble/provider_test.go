package ensemble

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func testParseUri(s string) string {
	hosts, path, err := parseZookeeperURI(s)
	So(err, ShouldBeNil)
	So(len(hosts), ShouldEqual, 3)
	return path
}

func TestFixedEnsembleProvider(t *testing.T) {

	Convey("A fixed ensemble provider should", t, func() {
		Convey("when parsing the uri string", func() {

			Convey("parse a complete uri string", func() {
				s := "zk://blah-blah-123:2181,foo-3939:2181,bar-bar-eei494:2181/a-path-here"
				path := testParseUri(s)
				So(path, ShouldEqual, "/a-path-here")
			})

			Convey("parse a uri with a zookeeper scheme", func() {
				s := "zookeeper://blah-blah-123:2181,foo-3939:2181,bar-bar-eei494:2181/a-path-here"
				path := testParseUri(s)
				So(path, ShouldEqual, "/a-path-here")
			})

			Convey("parse a uri with a missing scheme", func() {
				s := "blah-blah-123:2181,foo-3939:2181,bar-bar-eei494:2181/a-path-here"
				path := testParseUri(s)
				So(path, ShouldEqual, "/a-path-here")
			})

			Convey("parse a uri without a path", func() {
				s := "blah-blah-123:2181,foo-3939:2181,bar-bar-eei494:2181"
				path := testParseUri(s)
				So(path, ShouldEqual, "/")
			})

			Convey("parse a uri with single host", func() {
				s := "blah-blah-123:2181"
				hosts, path, err := parseZookeeperURI(s)
				So(err, ShouldBeNil)
				So(len(hosts), ShouldEqual, 1)
				So(path, ShouldEqual, "/")
			})

			Convey("fail to parse when there are no hosts", func() {
				s := "zookeeper://"
				_, _, err := parseZookeeperURI(s)
				So(err, ShouldNotBeNil)
			})
		})

		Convey("when creating", func() {
			s := "zk://blah-blah-123:2181,foo-3939:2181,bar-bar-eei494:2181/a-path-here"
			prov, pth, err := Fixed(s)
			So(err, ShouldBeNil)
			So(pth, ShouldEqual, "/a-path-here")

			Convey("start should not have an error", func() {
				e := prov.Start()
				So(e, ShouldBeNil)
			})

			Convey("stop should not have an error", func() {
				e := prov.Close()
				So(e, ShouldBeNil)
			})

			Convey("start should not have an error", func() {
				e := prov.Hosts()
				So(len(e), ShouldEqual, 3)
			})

		})
	})
}
