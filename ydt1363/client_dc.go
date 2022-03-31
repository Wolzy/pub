package gYDT1363


// 获取直流配电系统的浮点数
func (this *Client) GetDCFloat(addr int) (int, error) {
	return this.getData(addr, Cid1DC, Cid2GetFloat, "", 0)

}

// 获取直流配电系统的整形数
func (this *Client) GetDCFixed(addr int) (int, error) {
	return this.getData(addr, Cid1DC, Cid2GetFixed, "", 0)

}

// 获取直流配电系统的状态量
func (this *Client) GetDCState(addr int) (int, error) {
	return this.getData(addr, Cid1DC, Cid2GetState, "", 0)

}

// 获取直流配电系统的告警量
func (this *Client) GetDCWarn(addr int) (int, error) {
	return this.getData(addr, Cid1DC, Cid2GetWarn, "", 0)

}
