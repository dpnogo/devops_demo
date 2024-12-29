package mysql

import (
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"iam/internal/apiserver/store"
	"iam/internal/pkg/options"
	"sync"

	"iam/pkg/db"
)

// datastore 实现
type datastore struct {
	db *gorm.DB
}

func (store *datastore) User() store.UserStore {
	return newUsers(store.db)
}

func (store *datastore) Close() error {

	myDb, err := store.db.DB()
	if err != nil {
		return errors.Wrap(err, "get gorm db instance failed")
	}
	return myDb.Close()

}

var (
	mysqlFactory store.Factory
	once         = &sync.Once{}
)

// GetStoreFactory 根据配置获取 store 的Factory接口
func GetStoreFactory(opts *options.MysqlOptions) (store.Factory, error) {

	// 判断options是否为空，若为空则说明有其他地方已经执行过，判断
	if opts == nil && store.GetFactory() == nil {
		return nil, fmt.Errorf("mysql options and store factory is empty")
	}

	var (
		err    error
		gormDb *gorm.DB
	)

	// 首次根据 mysql option 进行初始化mysql client
	once.Do(func() {
		// 调用外部通用的 初始化 mysql 的库
		gormDb, err = db.NewDb(db.Options{
			Host:                  opts.Host,
			Username:              opts.Username,
			Password:              opts.Password,
			Database:              opts.Database,
			MaxIdleConnections:    opts.MaxIdleConnections,
			MaxOpenConnections:    opts.MaxOpenConnections,
			MaxConnectionLifeTime: opts.MaxConnectionLifeTime,
			LogLevel:              opts.LogLevel,
			// Logger:                logger.New(opts.LogLevel),  // 自己实现的logger
		})
		mysqlFactory = &datastore{gormDb}
	})

	if mysqlFactory == nil || err != nil {
		return nil, fmt.Errorf("failed to get mysql store fatory, mysqlFactory: %+v, error: %w", mysqlFactory, err)
	}

	return mysqlFactory, nil
}
