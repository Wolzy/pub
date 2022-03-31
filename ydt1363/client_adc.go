package gYDT1363


// 获取整流配电系统的浮点数
func (this *Client) GetADCFloat(addr int) (int, error) {
	return this.getData(1, Cid1ADC, Cid2GetFloat, "", 0)

}

// 获取整流配电系统的整形数
func (this *Client) GetADCFixed(addr int) (int, error) {
	return this.getData(1, Cid1ADC, Cid2GetFixed, "", 0)

}

// 获取整流配电系统的状态量
func (this *Client) GetADCState(addr int) (int, error) {
	return this.getData(1, Cid1ADC, Cid2GetState, "", 0)

}

// 获取整流配电系统的告警量
func (this *Client) GetADCWarn() (int, error) {
	return this.getData(1, Cid1ADC, Cid2GetWarn, "", 0)

}
