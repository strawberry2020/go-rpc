package rpc

import (
	"fmt"
	"net"
	"sync"
	"testing"
)

func TestSession_ReadWrite(t *testing.T) {
	// 定义监听IP和端口
	addr := "127.0.0.1:8000"
	// 定义传输的数据
	my_data := "hello"
	// 等待组
	wg := sync.WaitGroup{}
	// 携程 1个读 1个写
	wg.Add(2)
	// 写数据携程
	go func() {
		defer wg.Done()
		// 创建tcp连接
		lis, err := net.Listen("tcp", addr)
		if err != nil {
			t.Fatal(err)
		}
		conn, _ := lis.Accept()
		s := Session{conn: conn}
		// 写数据
		err = s.Write([]byte(my_data))
		if err != nil {
			t.Fatal(err)
		}
	}()

	// 读数据携程
	go func() {
		defer wg.Done()
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			t.Fatal(err)
		}
		s := Session{conn: conn}
		// 读数据
		data, err := s.Read()
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(string(data))
	}()
	wg.Wait()
}
