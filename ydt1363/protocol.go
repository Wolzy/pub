package gYDT1363

import (
	"errors"
	"fmt"
)

// YD/T 1363.3-2005
// SOI(1) + VER(1) + ADR(1) + CID1(1) + CID2(1) + LENGTH(2) + INFO(X) + CHKSUM(2) + EOI(1)
func NewProt1363(adr, cid1, cid2 string, length uint16, data string) *Prot1363 {
	return &Prot1363{
		CodeSOI:  0x7E, // ~
		CodeVER:  "20",
		CodeEOI:  0x0D, // 回车符
		CodeADR:  adr,
		CodeCID1: cid1,
		CodeCID2: cid2,
		CodeInfo: data,
		Length:   length,
	}
}

// 计算Length
// LENGTH共2个字节，由LENID和LCHKSUM组成，LENID表示INFO项的传送的ASCII码字节数
// 当LENID=0时，INFO为空，即无该项。LENGTH传输中先传高字节，再传低字节，分四个ASCII码传送。
// 校验码的计算：D11D10D9D8+D7D6D5D4+D3D2D1D0，求和后模16的余数取反加1
func (this *Prot1363) SumLength() error {
	return nil
}

// VER + ADR + CID1 + CID2 + LENGTH(
func (this *Prot1363) Chksum() uint16 {
	var tmp uint16
	tmp = uint16(this.CodeVER[0]) + uint16(this.CodeVER[1]) +
		uint16(this.CodeADR[0]) + uint16(this.CodeADR[1]) +
		uint16(this.CodeCID1[0]) + uint16(this.CodeCID1[1]) +
		uint16(this.CodeCID2[0]) + uint16(this.CodeCID2[1]) +
		uint16(this.CodeLen[0]) +
		uint16(this.CodeLen[1]) +
		uint16(this.CodeLen[2]) +
		uint16(this.CodeLen[3])

	for i := 0; i < len(this.CodeInfo); i++ {
		tmp += uint16(this.CodeInfo[i])
	}

	tmp = ^(tmp & 0xFFFF) + 1
	return tmp
}

func (this *Prot1363) Send() error {

	return nil
}

func (this *Prot1363) Serial() ([]byte, int) {
	var totalLen int
	var sendData = make([]byte, 1024)
	sendData[0] = this.CodeSOI
	totalLen = 1

	// 计算长度并转换成ASCII
	this.CodeLen = fmt.Sprintf("%04X", Lchksum(this.Length)*16*16*16+this.Length)

	for i := 0; i < len(this.CodeVER); i++ {
		sendData[totalLen+i] = this.CodeVER[i]
	}
	totalLen = totalLen + len(this.CodeVER)

	for i := 0; i < len(this.CodeADR); i++ {
		sendData[totalLen+i] = this.CodeADR[i]
	}
	totalLen = totalLen + len(this.CodeADR)

	for i := 0; i < len(this.CodeCID1); i++ {
		sendData[totalLen+i] = this.CodeCID1[i]
	}
	totalLen = totalLen + len(this.CodeCID1)

	for i := 0; i < len(this.CodeCID2); i++ {
		sendData[totalLen+i] = this.CodeCID2[i]
	}
	totalLen = totalLen + len(this.CodeCID2)

	for i := 0; i < len(this.CodeLen); i++ {
		sendData[totalLen+i] = this.CodeLen[i]
	}
	totalLen = totalLen + len(this.CodeLen)

	for i := 0; i < len(this.CodeInfo); i++ {
		sendData[totalLen+i] = this.CodeInfo[i]
	}
	totalLen = totalLen + len(this.CodeInfo)

	// 计算checksum
	this.CodeChk = fmt.Sprintf("%04X", this.Chksum())
	for i := 0; i < len(this.CodeChk); i++ {
		sendData[totalLen+i] = this.CodeChk[i]
	}
	totalLen = totalLen + len(this.CodeChk)

	sendData[totalLen] = this.CodeEOI
	totalLen++

	return sendData, totalLen
}

func Prot1363Error(code string) error {
	switch code {
	case "01":
		return errors.New("version error")
	case "02":
		return errors.New("chksum error")
	case "03":
		return errors.New("lchksum error")
	case "04":
		return errors.New("cid2 invalid")
	case "05":
		return errors.New("cmd format error")
	case "06":
		return errors.New("data error")
	case "E0":
		return errors.New("permission denied")
	case "E1":
		return errors.New("operation failed")
	case "E2":
		return errors.New("device error")
	case "E3":
		return errors.New("device write protected")
	default:
		return errors.New("reserved error code")
	}
}

func (this *Prot1363) String() string {
	return fmt.Sprintf("%02X%02X%02X%02X%02X%02X%02X%02X%02X",
		this.CodeSOI,
		this.CodeVER,
		this.CodeADR,
		this.CodeCID1,
		this.CodeCID2,
		this.CodeLen,
		this.CodeInfo,
		this.CodeChk,
		this.CodeEOI)
}