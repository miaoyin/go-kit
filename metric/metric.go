package metric


//Provider 监控接口
type Provider interface{
	Metrics() any
}
