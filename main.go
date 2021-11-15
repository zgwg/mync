package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

var serverConn *net.TCPConn
var isConnetion bool
func tcpServer(port string) {
	addr,err := net.ResolveTCPAddr("tcp",":"+port)
	if err!=nil{
		fmt.Println("监听地址错误！",err)
		os.Exit(-1)
		return
	}
	lis,err := net.ListenTCP("tcp",addr)
	if err!=nil{
		fmt.Println("打开监听错误！",err)
		os.Exit(-1)
		return
	}
	fmt.Println("已经监听端口"+port+"...")
	for {
		conn, err := lis.AcceptTCP()
		if err!=nil{
			fmt.Println("接收客服端错误！",err)
			os.Exit(-1)
			return
		}
		fmt.Println("已经连接客服端"+conn.RemoteAddr().String())
		isConnetion = true
		serverConn = conn
		for {
			buff := make([]byte, 1500)
			bufLen, err := conn.Read(buff)
			if err != nil {
				fmt.Println("已经断开客服端"+conn.RemoteAddr().String())
				isConnetion = false
				conn.Close()
				break
			}
			if buff[0]==3{
				fmt.Println("已经断开客服端"+conn.RemoteAddr().String())
				isConnetion = false
				conn.Close()
				break
			}
			buff = buff[:bufLen]
			fmt.Print(string(buff))
		}
	}

}
func readServer(conn *net.TCPConn){
	for {
		buff := make([]byte, 1500)
		bufLen, err := conn.Read(buff)
		if err != nil {
			os.Exit(-1)
			return
		}
		buff = buff[:bufLen]
		fmt.Println(string(buff))
	}
}
func scanAllTcpPort(ipAddr string)[]uint32{
	var port uint32
	var goNum int32

	var res []uint32
	var lock sync.Mutex
	var wg sync.WaitGroup
	wg.Add(65535)
	for port = 1;port<=65535;port++{
		if atomic.LoadInt32(&goNum)>10000{
			time.Sleep(time.Millisecond*10)
		}
		go func(p uint32){
			atomic.AddInt32(&goNum,1)
			defer wg.Done()
			defer atomic.AddInt32(&goNum,-1)
			conn,err := net.DialTimeout("tcp",ipAddr+":"+strconv.Itoa(int(p)),time.Second*2)
			if err==nil{
				lock.Lock()
				res = append(res,p)
				lock.Unlock()
				conn.Close()
			}
		}(port)
	}

	wg.Wait()
	return res
}
func main() {

	isListen := flag.Bool("l",false,"监听状态")
	listenPort :=flag.Int("p",9999,"监听端口号")
	isScan :=flag.Bool("s",false,"扫描所有打开的TCP端口")
	flag.Parse()
	if *isListen {
		go tcpServer(strconv.Itoa(*listenPort))

	} else if *isScan{
		args := flag.Args()
		argLen := len(args)
		if argLen == 0{
			fmt.Println("用法:客户端：mync 主机地址:端口号")
			fmt.Println("     服务器：mync -l -p 监听端口号")
			fmt.Println("     扫描所有TCP端口：mync -s 主机地址")
			return
		}
		res :=scanAllTcpPort(args[0])
		fmt.Println("打开的TCP端口号为：",res)
		return
	} else{
		args := flag.Args()
		argLen := len(args)
		if argLen == 0{
			fmt.Println("用法:客户端：mync 主机地址:端口号")
			fmt.Println("     服务器：mync -l -p 监听端口号")
			fmt.Println("     扫描所有TCP端口：mync -s 主机地址")
			return
		}
		_ ,err :=net.ResolveTCPAddr ("tcp",args[0])
		if err!=nil{
			fmt.Println("用法:客户端：mync 主机地址:端口号")
			fmt.Println("     服务器：mync -l -p 监听端口号")
			fmt.Println("     扫描所有TCP端口：mync -s 主机地址")
			return
		}

		ipconn,err :=net.DialTimeout("tcp",args[0],time.Second*5)

		//conn,err :=net.DialTCP("tcp",nil,addr)
		if err!=nil {
			fmt.Println("连接主机错误：",err)
			return
		}
		conn := ipconn.(*net.TCPConn)
		isConnetion = true
		fmt.Println("已经连接到服务器!")
		serverConn = conn
		go readServer(conn)

	}
	for {
		if serverConn==nil || !isConnetion{
			time.Sleep(100)
			continue
		}
		reader := bufio.NewReader(os.Stdin)
		res, _ := reader.ReadString('\n')
		_ ,err := serverConn.Write([]byte(res))
		if err!=nil{
			fmt.Println("已断开连接，发送不成功！")
		}
	}
}