package gYDT1363

import (
	"errors"
	"fmt"
	"github.com/tarm/serial"
	"io"
	"time"
)

type Client struct {
	Port     string
	BaudRate int
	TimeOut  time.Duration
	iorwc    io.ReadWriteCloser
	CmdSend  *Prot1363
	CmdRecv  *Prot1363
}

func NewClient(port string, rate int, timeout time.Duration) Client {
	return Client{
		Port:     port,
		BaudRate: rate,
		TimeOut:  timeout,
		//Prot: NewProt1363()
	}
}

func (this *Client) Connect() error {
	var err error
	comCfg := &serial.Config{Name: this.Port, Baud: this.BaudRate, ReadTimeout: this.TimeOut * time.Second}
	this.iorwc, err = serial.OpenPort(comCfg)
	if err != nil {
		return err
	}

	return nil
}

func (this *Client) Close() error {
	if err := this.iorwc.Close(); err != nil {
		return err
	}

	return nil
}

func (this *Client) write(buf []byte) (int ,error) {
	return this.iorwc.Write(buf)
}

func (this *Client) read() (int, error) {
	var cnt, ret int
	var err error

	// 除了SOI是16进制，其他都是BCD码，一个字节数据占两个字节
	// SOI(0x7E), VER(0x032,0x30), ADR(0x30,0x31), CID1(0x34,0x??), CID2(0x??,0x??)
	// CID2为返回值
	var head = make([]byte, 9)
	cnt, err = this.iorwc.Read(head)
	if err != nil {
		fmt.Println("读取回执文件头失败")
		fmt.Println("head:", head)
		return -1, err
	}
	fmt.Printf("Recv Head: %X\n", head)

	// SOI不匹配
	if head[0] != 0x7E {
		return 0, errors.New("recv soi error")
	}
	this.CmdRecv.CodeSOI = head[0]
	// VER不匹配
	if head[1] != this.CmdSend.CodeVER[0] || head[2] != this.CmdSend.CodeVER[1] {
		return 0, errors.New("recv ver error")
	}
	this.CmdRecv.CodeVER = string(head[1:3])
	// ADR不匹配
	if head[3] != this.CmdSend.CodeADR[0] || head[4] != this.CmdSend.CodeADR[1] {
		return 0, errors.New("recv adr error")
	}
	this.CmdRecv.CodeADR = string(head[3:5])
	// CID1不匹配
	if head[5] != this.CmdSend.CodeCID1[0] || head[6] != this.CmdSend.CodeCID1[1] {
		return 0, errors.New("recv cid1 error")
	}
	this.CmdRecv.CodeCID1 = string(head[5:7])
	// CID2不匹配
	if head[7] != 0x30 && head[8] != 0x30 {
		var tmpCID2 = string(head[7:9])
		err = Prot1363Error(tmpCID2)
		return 0, err
	}
	this.CmdRecv.CodeCID2 = string(head[7:9])
	ret += cnt

	// LENGTH(4字节)
	var length = make([]byte, 4)
	cnt, err = this.iorwc.Read(length)
	fmt.Println("length: ", length)

	if err != nil {
		fmt.Println("读取回执信息长度失败")
		return -1, err
	}
	lchksum := length[0] - 0x30
	this.CmdRecv.Length = uint16(length[1]-0x30)*16*16 + uint16(length[2]-0x30)*16 + uint16(length[3]-0x30)
	fmt.Printf("lchksum: %02X;%d\n", lchksum, this.CmdRecv.Length)
	if lchksum != byte(Lchksum(this.CmdRecv.Length)) {
		return 0, errors.New("recv lchksum error")
	}
	this.CmdRecv.CodeLen = string(length)
	fmt.Printf("Length: %d\n", this.CmdRecv.Length)

	// 字节头13个
	ret += cnt

	// 数据域长度+4字节chksum+1字节EOI(0x0D)
	cnt = 0
	var data = make([]byte, 0)
	var buf = make([]byte, 32)
	for {
		num, err := this.iorwc.Read(buf)
		if err != nil {
			return 0, err
		}
		if num == 0 {
			break
		}
		fmt.Printf("%02X\n", buf)
		data = append(data, buf...)
		fmt.Printf("%02X\n", data)
		cnt += num
	}

	fmt.Printf("Data: %02X\n", data)
	if data[len(buf)-1] != 0x0D {
		return 0, errors.New("recv EOI error")
	}
	this.CmdRecv.CodeInfo = string(data[:len(data)-5])
	this.CmdRecv.CodeChk = string(data[len(data)-5:len(data)-1])
	this.CmdRecv.CodeEOI = data[len(data)-1]
	ret += len(data)

	return ret, nil
}

func (this *Client) getData(Addr, Cid1, Cid2 int, Data string, Length uint16) (int, error) {
	addr := fmt.Sprintf("%02X",Addr)
	cid1 := fmt.Sprintf("%02X",Cid1)
	cid2 := fmt.Sprintf("%02X",Cid2)
	this.CmdSend = nil
	this.CmdRecv = nil
	this.CmdSend = NewProt1363(addr, cid1, cid2, Length, Data)
	this.CmdRecv = NewProt1363(addr, cid1, "00",0, "")
	pbuf, plen := this.CmdSend.Serial()
	wbuf := pbuf[:plen]

	ret, err := this.write(pbuf)
	if err != nil {
		return ret, err
	}
	fmt.Printf("Cmd Send: %X\n", wbuf)

	// Wait 0.5s
	time.Sleep(500 * time.Millisecond)

	// Read response
	ret, err = this.read()
	if err != nil {
		return ret, err
	}
	fmt.Printf("%s\n", this.CmdRecv.String())


	return ret, err
}
