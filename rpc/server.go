package rpc

import (
	"net"
	"reflect"
)

// 声明服务端
type Server struct {
	// 地址
	addr string
	// 服务端维护函数名到函数反射值的map
	funcs map[string]reflect.Value
}

// 创建服务端对象
func NewServer(addr string) *Server {
	return &Server{addr: addr, funcs: make(map[string]reflect.Value)}
}

// 服务端绑定注册方法
// 将函数名与函数真正实现对应起来
// 第一个参数为函数名，第二个传入真正的函数
func (s *Server) Register(rpcName string, f interface{}) {
	if _, ok := s.funcs[rpcName]; ok {
		return
	}
	// map中没有值，则将映射添加到map，便于调用
	fVal := reflect.ValueOf(f)
	s.funcs[rpcName] = fVal
}

// 服务端等待调用
func (s *Server) Run() {
	// 监听
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return
	}
	for {
		// 获取连接
		conn, err := lis.Accept()
		if err != nil {
			return
		}
		// 创建会话
		srvSession := NewSession(conn)
		// PRC读数据
		b, err := srvSession.Read()
		if err != nil {
			return
		}
		// 对数据解码
		rpcData, err := decode(b)
		// 通过Name获取函数名
		f, ok := s.funcs[rpcData.Name]
		if !ok {
			return
		}
		// 解析遍历客户端传的参数，放到数组中
		inArgs := make([]reflect.Value, 0, len(rpcData.Args))
		for _, arg := range rpcData.Args {
			inArgs = append(inArgs, reflect.ValueOf(arg))
		}
		// 反射调用方法，传入参数
		out := f.Call(inArgs)
		// 解析遍历执行结果，保存到数组中
		outArgs := make([]interface{}, 0, len(out))
		for _, o := range out {
			outArgs = append(outArgs, o.Interface())
		}
		// 包装数据，返回给客户端
		resRPCData := RPCData{rpcData.Name, outArgs}
		// 编码
		respBytes, err := encode(resRPCData)
		if err != nil {
			return
		}
		// 使用RPC写数据
		err = srvSession.Write(respBytes)
		if err != nil {
			return
		}
	}
}
