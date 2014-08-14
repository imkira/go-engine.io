package engineio

import (
	"io"
	"net/http"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

type mockWriteCloser struct {
}

func (wc *mockWriteCloser) Write(p []byte) (int, error) {
	return len(p), nil
}

func (wc *mockWriteCloser) Close() error {
	return nil
}

type mockTransport struct {
	conn   Conn
	writer io.WriteCloser
}

func (t *mockTransport) Name() string {
	return "mock"
}

func (t *mockTransport) SetConn(conn Conn) {
	t.conn = conn
}

func (t *mockTransport) ServeHTTP(http.ResponseWriter, *http.Request) {
}

func (t *mockTransport) Close() error {
	t.conn = nil
	return nil
}

func (t *mockTransport) NextWriter(msgType MessageType, packetType packetType) (io.WriteCloser, error) {
	return t.writer, nil
}

func TestConn(t *testing.T) {
	Convey("Packet handling", t, func() {
		srv, err := NewServer(nil)
		So(err, ShouldBeNil)
		So(srv, ShouldNotBeNil)

		fwc := &mockWriteCloser{}
		ft := &mockTransport{writer: fwc}
		So(ft, ShouldNotBeNil)

		Convey("OPEN", func() {
			srv.config.PingTimeout = 1 * time.Millisecond
			c, err := newConn("1", srv, ft, nil)
			So(err, ShouldBeNil)
			So(c, ShouldNotBeNil)
			dec := &packetDecoder{}
			dec.t = _OPEN
			c.onPacket(dec)
		})

		Convey("PONG", func() {
			srv.config.PingTimeout = 1 * time.Millisecond
			c, err := newConn("1", srv, ft, nil)
			So(err, ShouldBeNil)
			So(c, ShouldNotBeNil)
			dec := &packetDecoder{}
			dec.t = _PONG
			c.onPacket(dec)
		})
	})
}
