package gYDT1363

func Lchksum(len uint16) uint16 {
	// 高字节低4位(/256)+低字节高4位(/16)+低字节低4位(%16)
	s1 := len/256 + len/16 + len%16
	// 模16余数取反加1
	s2 := ^(s1 % 16) + 1
	// 防止出现0x100，模16
	s3 := s2 % 16

	return s3
}
