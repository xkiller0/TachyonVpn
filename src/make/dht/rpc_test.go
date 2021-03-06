package dht

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"net"
	"testing"
)

func TestRpcNodeFindNode_one_to_one(t *testing.T) {
	node1 := newPeerNode(newPeerNodeRequest{
		id:   1,
		port: 60001,
		bootstrapRpcNodeList: []*rpcNode{
			{
				Id:   2,
				Ip:   net.ParseIP("127.0.0.1"),
				Port: 60002,
			},
			{
				Id:   4,
				Ip:   net.ParseIP("127.0.0.1"),
				Port: 60004,
			},
		},
	})
	close1 := node1.StartRpcServer()
	defer close1()
	node4 := newPeerNode(newPeerNodeRequest{
		id:   4,
		port: 60004,
		bootstrapRpcNodeList: []*rpcNode{
			{
				Id:   2,
				Ip:   net.ParseIP("127.0.0.1"),
				Port: 60002,
			},
		},
	})
	close4 := node4.StartRpcServer()
	defer close4()
	//noinspection SpellCheckingInspection
	data := []byte("1drnk7yc53frym6qe2saupptppytj7cbk")
	dataKey := hash(data)
	node1.store(data)
	node3 := newPeerNode(newPeerNodeRequest{
		id: 3,
		bootstrapRpcNodeList: []*rpcNode{
			{
				Id:   1,
				Ip:   net.ParseIP("127.0.0.1").To4(),
				Port: 60001,
			},
		},
	})
	closestRpcNodeList := node3.findNode(2)
	udwTest.Equal(len(closestRpcNodeList), 2)
	udwTest.Equal(closestRpcNodeList[0].Id, uint64(2))
	udwTest.Equal(closestRpcNodeList[0].Port, uint16(60002))
	v := node3.findValue(dataKey)
	udwTest.Equal(string(v), string(data))
}

func TestRpcServerRandomPort(t *testing.T) {
	node1 := newPeerNode(newPeerNodeRequest{
		id:   1,
	})
	close1 := node1.StartRpcServer()
	defer close1()
	udwTest.Ok(node1.port!=0)
}

//var responseTimeoutError = errors.New("timeout")

//func debugClientSend(request []byte, afterWrite func(conn net.Conn) (isReturn bool)) (response []byte, err error) {
//	conn, err := net.Dial("udp", "127.0.0.1:"+strconv.Itoa(rpcPort))
//	udwErr.PanicIfError(err)
//	_, err = conn.Write(request)
//	udwErr.PanicIfError(err)
//	if afterWrite != nil {
//		isReturn := afterWrite(conn)
//		if isReturn {
//			return
//		}
//	}
//	buf := make([]byte, 2<<10)
//	err = conn.SetDeadline(time.Now().Add(time.Millisecond * 300))
//	udwErr.PanicIfError(err)
//	n, err := conn.Read(buf)
//	if err != nil {
//		return nil, responseTimeoutError
//	}
//	return buf[:n], nil
//}
//
//func TestRpcNodeErrorClient(t *testing.T) {
//	node := newPeerNode(0)
//	closeRpcServer := node.StartRpcServer()
//	defer closeRpcServer()
//	errMsg := ""
//	_, err := debugClientSend([]byte("1"), nil)
//	if err != nil {
//		errMsg = err.Error()
//	}
//	udwTest.Equal(errMsg, responseTimeoutError.Error())
//	_, err = debugClientSend([]byte{0x02, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, nil)
//	if err != nil {
//		errMsg = err.Error()
//	}
//	udwTest.Equal(errMsg, responseTimeoutError.Error())
//	_, err = debugClientSend([]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, nil)
//	if err != nil {
//		errMsg = err.Error()
//	}
//	udwTest.Equal(errMsg, responseTimeoutError.Error())
//	_, err = debugClientSend([]byte{0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01}, func(conn net.Conn) bool {
//		_ = conn.Close()
//		return true
//	})
//	udwTest.Equal(err, nil)
//}
//
//func debugServerRespond(correctIdMessage bool, response []byte) (close func()) {
//	closer := udwClose.NewCloser()
//	packetConn, err := net.ListenPacket("udp", ":"+strconv.Itoa(rpcPort))
//	udwErr.PanicIfError(err)
//	closer.AddOnClose(func() {
//		_ = packetConn.Close()
//	})
//	go func() {
//		rBuf := make([]byte, 2<<10)
//		n, addr, err := packetConn.ReadFrom(rBuf)
//		udwErr.PanicIfError(err)
//		request := rpcMessage{}
//		err = request.rpcMessageDecode(rBuf[:n])
//		udwErr.PanicIfError(err)
//		if correctIdMessage && len(response) > 5 {
//			binary.BigEndian.PutUint32(response[1:5], request._idMessage)
//		}
//		_, err = packetConn.WriteTo(response, addr)
//		udwErr.PanicIfError(err)
//	}()
//	return closer.Close
//}
//
//func TestRpcNodeErrorServer(t *testing.T) {
//	rNode2 := rpcNode{
//		id: 1,
//		ip: "127.0.0.1",
//	}
//	errMsg := ""
//	_, err := rNode2.findNode(2)
//	if err != nil {
//		errMsg = err.Error()
//	}
//	udwTest.Ok(strings.Contains(errMsg, errorRpcCallResponseTimeout))
//
//	_close := debugServerRespond(false, []byte("1"))
//	_, err = rNode2.findNode(2)
//	if err != nil {
//		errMsg = err.Error()
//	}
//	udwTest.Ok(strings.Contains(errMsg, errorRpcCallResponseTimeout))
//	_close()
//
//	_close = debugServerRespond(true, []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff})
//	err = rNode2.store([]byte("123"))
//	if err != nil {
//		errMsg = err.Error()
//	}
//	udwTest.Ok(strings.Contains(errMsg, errorRpcCallResponseTimeout))
//	_close()
//
//	_close = debugServerRespond(false, []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff})
//	_, _, err = rNode2.findValue(2)
//	if err != nil {
//		errMsg = err.Error()
//	}
//	udwTest.Ok(strings.Contains(errMsg, errorRpcCallResponseTimeout))
//	_close()
//
//	_close = debugServerRespond(true, []byte{cmdOk, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff})
//	_, err = rNode2.findNode(2)
//	if err != nil {
//		errMsg = err.Error()
//	}
//	udwTest.Ok(strings.Contains(errMsg, "fhf1b2xk9u9"))
//	_close()
//
//	_close = debugServerRespond(true, []byte{cmdOk, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff})
//	_, _, err = rNode2.findValue(2)
//	if err != nil {
//		errMsg = err.Error()
//	}
//	udwTest.Ok(strings.Contains(errMsg, "kge9ma4b69"))
//	_close()
//}
