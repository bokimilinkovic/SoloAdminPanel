package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

// CreateDBConnection creates and loads the gorm object with the given configuration.
// Driver has to be imported first from the client code (e.g. import _ "github.com/go-sql-driver/mysql").
func CreateDBConnection(config *Config, logger logrus.FieldLogger) (*gorm.DB, error) {
	db, err := gorm.Open("postgres", config.ConnectionURL())
	if err != nil {
		return nil, err
	}

	db.LogMode(config.LogMode)
	db.SetLogger(&gormLogger{logger})

	return db, nil
}

// gormLogger is a Logrus logger that implements the gorm interface for customized query logging.
// Source: https://gist.github.com/bnadland/2e4287b801a47dcfcc94.
type gormLogger struct {
	logrus.FieldLogger
}

// Print implements the gorm.LogWriter interface adjusted for query logging.
// It adds entry in log like this:
// `{"done_in":2207375,"level":"info","msg":"SELECT count(*) FROM `clients` WHERE (name = ?)","params":["joe"],...
// ..."source":"/path/file.go:65","time":"2019-02-10T18:42:21Z","type":"sql"}`.
func (g *gormLogger) Print(v ...interface{}) {
	if v[0] == "sql" {
		g.WithFields(logrus.Fields{"type": v[0], "source": v[1], "done_in": v[2], "params": v[4]}).Info(v[3])
	} else if v[0] == "log" {
		g.WithFields(logrus.Fields{"type": v[0], "source": v[1]}).Info(v[2])
	}
}
