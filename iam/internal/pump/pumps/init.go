package pumps

var pumpsTypes map[string]Pump // 根据情况进行选择对应的下游组件

func init() {
	pumpsTypes = make(map[string]Pump)

	// 根据需要添加组件

	//pumpsTypes["elasticsearch"] = &ElasticsearchPump{}
	//pumpsTypes["kafka"] = &KafkaPump{}

}
