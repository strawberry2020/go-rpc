package rpc

import (
	"net"
	"reflect"
)

// 声明客户端
type Client struct {
	conn net.Conn
}

// 创建客户端对象
func NewClient(conn net.Conn) *Client {
	return &Client{conn: conn}
}

// 实现通用的RPC客户端
// 绑定PRC访问的方法
// 传入访问的函数名

// 函数的具体实现在Server端，Client只有函数原型
// 使用MakeFunc() 完成原型到函数的调用

// fPtr指向函数原型
// xxx.callPRC("queryUser",&query)
func (c *Client) callPRC(rpcName string, fPtr interface{}) {
	// 通过反射，获取fPtr未初始化的函数原型
	fn := reflect.ValueOf(fPtr).Elem()
	// 另一个函数，作用是对第一个函数参数操作
	// 完成与Server交互
	f := func(args []reflect.Value) []reflect.Value {
		// 处理输入的参数
		inArgs := make([]interface{}, 0, len(args))
		for _, arg := range args {
			inArgs = append(inArgs, arg.Interface())
		}
		// 创建连接
		clientSession := NewSession(c.conn)
		// 编码数据
		reqRPC := RPCData{Name: rpcName, Args: inArgs}
		b, err := encode(reqRPC)
		if err != nil {
			panic(err)
		}
		// 写数据
		err = clientSession.Write(b)
		if err != nil {
			panic(err)
		}
		// 读取响应数据
		respBytes, err := clientSession.Read()
		if err != nil {
			panic(err)
		}
		// 解码数据
		respRPC, err := decode(respBytes)
		if err != nil {
			panic(err)
		}
		// 处理服务端返回的数据
		outArgs := make([]reflect.Value, 0, len(respRPC.Args))
		for i, arg := range respRPC.Args {
			// 必须进行nil转换
			if arg == nil {
				// 必须填充一个真正的类型，不能是nil
				outArgs = append(outArgs, reflect.Zero(fn.Type().Out(i)))
				continue
			}
			outArgs = append(outArgs, reflect.ValueOf(arg))
		}
		return outArgs
	}

	// 参数1:一个未初始化函数的方法值，类型是reflect.Type
	// 参数2:另一个函数，作用是对第一个函数参数操作
	// 返回 reflect.Value 类型
	// MakeFunc 使用传入的函数原型，创建一个绑定参数2的新函数
	v := reflect.MakeFunc(fn.Type(), f)
	// 为函数fPtr赋值
	fn.Set(v)

}
