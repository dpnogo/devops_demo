package pumps

import "context"

// 某一种存储方式

type ElasticsearchOperator interface {
	processData(ctx context.Context, data []interface{}, esConf *ElasticsearchConf) error
}

type ElasticsearchPump struct {
	CommonPumpConfig                       // 选择配置
	esConf           ElasticsearchConf     // es config
	operator         ElasticsearchOperator // 抽象批量往es中写入数据接口
}

type ElasticsearchConf struct {
	BulkConfig       ElasticsearchBulkConfig `mapstructure:"bulk_config"`
	IndexName        string                  `mapstructure:"index_name"`
	ElasticsearchURL string                  `mapstructure:"elasticsearch_url"`
	DocumentType     string                  `mapstructure:"document_type"`
	AuthAPIKeyID     string                  `mapstructure:"auth_api_key_id"`
	AuthAPIKey       string                  `mapstructure:"auth_api_key"`
	Username         string                  `mapstructure:"auth_basic_username"`
	Password         string                  `mapstructure:"auth_basic_password"`
	EnableSniffing   bool                    `mapstructure:"use_sniffing"`
	RollingIndex     bool                    `mapstructure:"rolling_index"`
	DisableBulk      bool                    `mapstructure:"disable_bulk"`
}

type ElasticsearchBulkConfig struct {
	Workers       int `mapstructure:"workers"`
	FlushInterval int `mapstructure:"flush_interval"`
	BulkActions   int `mapstructure:"bulk_actions"`
	BulkSize      int `mapstructure:"bulk_size"`
}

func (pump *ElasticsearchPump) New() Pump {

	//	return ElasticsearchPump{}
	return nil
}

func (pump *ElasticsearchPump) Init(cfg interface{}) error {
	return nil
}

func (pump *ElasticsearchPump) WriteData(datas []any) error {
	// 从redis中获取 --> 写入

	return nil
}

/*
	New() Pump
	Init(interface{}) error      //
	WriteData(datas []any) error // 往下游系统写入数据
*/
