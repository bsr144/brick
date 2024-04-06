package config

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type (
	App struct {
		Name    string
		Version string
		Port    int
	}

	Database struct {
		Name     string
		Host     string
		Port     string
		Username string
		Password string
	}

	REST struct {
		Port     int
		BodySize int64
		Debug    bool
		Timeout  struct {
			Read  int64
			Write int64
			Idle  int64
		}
	}
	Logging struct {
		Level  int
		Format string
	}

	Driver struct {
		App      App
		Database Database
		Logging  Logging
	}

	Bootstrap struct {
		RESTServer *gin.Engine
		Log        *logrus.Logger
		SqlDB      *sql.DB
		Driver     *Driver
	}
)
