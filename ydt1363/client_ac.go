package gYDT1363


// 获取交流配电系统的浮点数
func (this *Client) GetACFloat(addr int) (int, error) {
	return this.getData(addr, Cid1AC, Cid2GetFloat, "FF", 2)
}

// 获取交流配电系统的整形数
func (this *Client) GetACFixed(addr int) (int, error) {
	return this.getData(addr, Cid1AC, Cid2GetFixed, "FF", 2)
}

// 获取交流配电系统的状态量
func (this *Client) GetACState(addr int) (int, error) {
	return this.getData(addr, Cid1AC, Cid2GetState, "FF", 2)
}

// 获取交流配电系统的告警量
func (this *Client) GetACWarn(addr int) (int, error) {
	return this.getData(addr, Cid1AC, Cid2GetWarn, "FF", 2)
}

// 获取设备版本
func (this *Client) GetVersion(addr int) (int, error) {
	return this.getData(addr, Cid1AC, Cid2GetVer, "", 0)
}

// 获取设备地址
func (this *Client) GetAddress(addr int) (int, error) {
	return this.getData(addr, Cid1AC, Cid2GetAddr, "", 0)
}
