package database

import (
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

const (
	DriverMySQL     = "mysql"
	DriverPostgres  = "postgres"
	DriverSnowflake = "snowflake"
)

// DB object
var (
	Master *DB
	Slave  *DB
)

type (
	DSNConfig struct {
		DSN string
	}

	DBConfig struct {
		SlaveDSN        string `json:"slave_dsn" mapstructure:"slave_dsn"`
		MasterDSN       string `json:"master_dsn" mapstructure:"master_dsn"`
		RetryInterval   int    `json:"retry_interval" mapstructure:"retry_interval"`
		MaxIdleConn     int    `json:"max_idle" mapstructure:"max_idle"`
		MaxConn         int    `json:"max_conn" mapstructure:"max_conn"`
		ConnMaxLifetime string `json:"conn_max_lifetime" mapstructure:"conn_max_lifetime"`
	}

	DB struct {
		DBConnection    *sqlx.DB
		DBString        string
		RetryInterval   int
		MaxIdleConn     int
		MaxConn         int
		ConnMaxLifetime time.Duration
		doneChannel     chan bool
	}

	Store struct {
		Master *sqlx.DB
		Slave  *sqlx.DB
	}
	Options struct {
		dbTx *sqlx.Tx
	}
)

func (s *Store) GetMaster() *sqlx.DB {
	return s.Master
}

func (s *Store) GetSlave() *sqlx.DB {
	return s.Slave
}

func New(cfg DBConfig, dbDriver string) *Store {
	masterDSN := cfg.MasterDSN

	var connMaxLifetime time.Duration
	if cfg.ConnMaxLifetime != "" {
		duration, err := time.ParseDuration(cfg.ConnMaxLifetime)
		if err != nil {
			log.Fatal("Invalid ConnMaxLifetime value: " + err.Error())
			return &Store{}
		}

		connMaxLifetime = duration
	}

	Master = &DB{
		DBString:        masterDSN,
		RetryInterval:   cfg.RetryInterval,
		MaxIdleConn:     cfg.MaxIdleConn,
		MaxConn:         cfg.MaxConn,
		ConnMaxLifetime: connMaxLifetime,
		doneChannel:     make(chan bool),
	}

	err := Master.ConnectAndMonitor(dbDriver)
	if err != nil {
		log.Fatal("Could not initiate Master DB connection: " + err.Error())
		return &Store{}
	}

	slaveDSN := cfg.SlaveDSN
	if len(slaveDSN) > 0 {
		Slave = &DB{
			DBString:        slaveDSN,
			RetryInterval:   cfg.RetryInterval,
			MaxIdleConn:     cfg.MaxIdleConn,
			MaxConn:         cfg.MaxConn,
			ConnMaxLifetime: connMaxLifetime,
			doneChannel:     make(chan bool),
		}
		err = Slave.ConnectAndMonitor(dbDriver)
		if err != nil {
			log.Fatal("Could not initiate Slave DB connection: " + err.Error())
			return &Store{}
		}

		return &Store{Master: Master.DBConnection, Slave: Slave.DBConnection}
	}

	return &Store{Master: Master.DBConnection}
}

// Connect to database
func (d *DB) Connect(driver string) error {
	db, err := sqlx.Open(driver, d.DBString)
	if err != nil {
		return fmt.Errorf("failed to open DB connection: %w", err)
	}

	db.SetMaxOpenConns(d.MaxConn)
	db.SetMaxIdleConns(d.MaxIdleConn)

	if d.ConnMaxLifetime > 0 {
		db.SetConnMaxLifetime(d.ConnMaxLifetime)
	}

	d.DBConnection = db

	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping DB: %w", err)
	}

	return nil
}

// Connect and monitor database
func (d *DB) ConnectAndMonitor(driver string) error {
	err := d.Connect(driver)
	if err != nil {
		log.Printf("Not connected to database %s, trying", d.DBString)
		return err
	}

	ticker := time.NewTicker(time.Duration(d.RetryInterval) * time.Second)
	go func() error {
		for {
			select {
			case <-ticker.C:
				if d.DBConnection == nil {
					d.Connect(driver)
				} else {
					err := d.DBConnection.Ping()
					if err != nil {
						log.Printf("[Error]: DB reconnect error: %v", err.Error())
						return err
					}
				}
			case <-d.doneChannel:
				return nil
			}
		}
	}()
	return nil
}
