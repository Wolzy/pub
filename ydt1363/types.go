package gYDT1363

const (
	Success = iota
	ErrVersion
	ErrCheSum
	ErrLChkSum
	ErrCID2
	ErrCmd
	ErrInvalidData
	ErrInvalidAuth = 0xE0
	ErrOperation
	ErrDeviceFault
	ErrWriteProtect
)

const (
	Cid1AC  = 0x40
	Cid1ADC = 0x41
	Cid1DC  = 0x42
	Cid1SUN = 0x43
	Cid1Env = 0x80
	Cid1EXT = 0xE1
	Cid1ATS = 0xE2
)

const (
	Cid2GetFloat = 0x41
	Cid2GetFixed = 0x42
	Cid2GetState = 0x43
	Cid2GetWarn  = 0x44
	Cid2GetCtrl  = 0x45
	Cid2GetVer   = 0x4F
	Cid2GetAddr  = 0x50
)

type Prot1363 struct {
	CodeSOI  byte   // 起始标志位（START OF INFORMATION），取值7EH
	CodeVER  string // 通讯协议版本号（2.1版），取值21H
	CodeADR  string // 设备地址描述（1-254，0、255保留）
	CodeCID1 string // 控制标识码（设备类型描述）
	CodeCID2 string // 命令信息：控制标识码（数据动作类型描述）
	CodeRTN  string // 响应信息：返回码RTN，跟CodeCID2位置相同
	CodeLen  string // INFO字节长度（包括LENID和LCHKSUM），LENID=0时为0000
	CodeInfo string // 命令或者数据
	CodeChk  string // 和校验码
	CodeEOI  byte   // 回车
	Length   uint16
}

type CmdInfo struct {
	CID1     string `json:"CID1"`
	CID2     string `json:"CID2"`
	DataInfo string `json:"DataInfo"`
}

type DataInfo struct {
	DataI     int64
	DataF     float64
	RunState  int16
	WarnState int16
}
