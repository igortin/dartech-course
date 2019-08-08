package darproxy

type Count struct {
	scorer int
}

func (c *Count) ResetCount() {
	c.scorer = 0
}

func (c *Count) IncreaseCount(){
	c.scorer++
}