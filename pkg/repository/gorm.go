package repository

import (
	"github.com/STRRL/poc-reference-to-grafana/pkg/entity"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)
import "gorm.io/driver/mysql"

type Gorm struct {
}

func ConnectToMysql() (*gorm.DB, error) {
	dsn := "root:root@tcp(localhost:13306)/POC_REF_TO_GRAFANA?charset=utf8mb4"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, errors.Wrapf(err, "connect to mysql")
	}

	err = db.AutoMigrate(
		&entity.GrafanaPanel{},
		&entity.GrafanaPanelVariable{},
		&entity.GrafanaPanelBinding{},
		&entity.GrafanaPanelBindingVariableWithValue{},
	)
	if err != nil {
		return nil, errors.Wrapf(err, "gorm auto migrate table schema")
	}

	return db, nil
}
