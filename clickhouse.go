package main

import (
	"fmt"

	"database/sql"

	"time"

	"github.com/kshvakov/clickhouse"
)

func connectClickDB() {
	var err error
	clickHouseDB, err = sql.Open("clickhouse", config.Clickhouse.ConnString)
	if err != nil {
		Logger.Fatalln("Clickhouse connection error: ", err)
	}

	Logger.Println("Connected to Clickhouse, pinging")
	if err := clickHouseDB.Ping(); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			Logger.Fatalf("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		} else {
			Logger.Fatalln("Clickhouse ping error: ", err)
		}
	}
	Logger.Println("Ping to Clickhouse - OK")

	createClickDB()
	createStatsTable()
}

func createClickDB() {
	_, err := clickHouseDB.Exec(`CREATE DATABASE IF NOT EXISTS ` + config.Clickhouse.DBName)
	if err != nil {
		Logger.Fatalln("Creating database error:", err)
	}
	_, err = clickHouseDB.Exec(`USE ` + config.Clickhouse.DBName)
	if err != nil {
		Logger.Fatalln("Selecting db for session error:", err)
	}
}

func getTablesList() []string {
	rows, err := clickHouseDB.Query("SHOW TABLES")
	if err != nil {
		Logger.Fatalln("Error on SHOW TABLES:", err)
	}
	var tables []string

	for rows.Next() {
		var tableRow string
		if err := rows.Scan(&tableRow); err != nil {
			Logger.Fatalln("Error on scanning SHOW TABLES:", err)
		}
		tables = append(tables, tableRow)
	}
	return tables
}

func createStatsTable() {
	_, err := clickHouseDB.Exec(`
CREATE TABLE IF NOT EXISTS stats (
  event_time            DateTime,
  event_date            Date,
  node                  String,
  status_code           Int32,
  used_time				Float64,
  host					String,
  uri					String,
  user_agent			String,
  client_ip				String
) ENGINE = MergeTree(event_date, (node, status_code, used_time, uri, event_time), 8192)
`)
	if err != nil {
		Logger.Fatalln("Creating Account table error:", err)
	}
}

func insertStatsToCH(node string, status_code int, event_time time.Time, used_time time.Duration, host string, uri string, client_ip string, user_agent string) {
	var tx, err = clickHouseDB.Begin()
	if err != nil {
		Logger.Fatalln("clickHouseDB.Begin error:", err)
	}

	query := fmt.Sprintf("INSERT INTO stats (" +
		"node, " +
		"status_code, " +
		"event_time, " +
		"event_date, " +
		"used_time, " +
		"host," +
		"uri," +
		"user_agent," +
		"client_ip" +
		") VALUES (?, ?, ?, ?, ?, ?, ?, ?)")

	stmt, err := tx.Prepare(query)

	if err != nil {
		Logger.Fatalln("err stmt:", err)
	}

	_, err = stmt.Exec(
		node,
		status_code,
		event_time,
		event_time,
		used_time.Seconds(),
		host,
		uri,
		user_agent,
		client_ip,
	)
	if err != nil {
		Logger.Fatalln("Error on Exec Statement:", err)
	}

	if err := tx.Commit(); err != nil {
		Logger.Fatalln("Error on Commit Statement:", err)
	}
}
