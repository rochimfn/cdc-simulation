package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

type PersonField struct {
	PersonId       uint   `json:"person_id"`
	FullName       string `json:"fullname"`
	NumTransaction uint   `json:"num_transaction"`
	CompanyId      uint   `json:"company_id"`
}

type DebeziumValueEvent struct {
	Payload struct {
		Before    *PersonField `json:"before"`
		After     *PersonField `json:"after"`
		Operation string       `json:"op"`
	} `json:"payload"`
}

type DebeziumKeyEvent struct {
	Payload struct {
		PersonId uint `json:"person_id"`
	} `json:"payload"`
}

type PersonTarget struct {
	PersonId       uint   `json:"person_id"`
	FullName       string `json:"fullname"`
	NumTransaction uint   `json:"num_transaction"`
	CompanyName    string `json:"company_name"`
}

func getCompany(db *sql.DB, companyId uint) (*string, error) {
	var companyName string
	q := db.QueryRow("select fullname from company where company_id=?", companyId)
	if q.Err() != nil {
		return nil, q.Err()
	} else {
		q.Scan(&companyName)
		return &companyName, nil
	}
}

func insertPerson(db *sql.DB, p PersonTarget) error {
	_, err := db.Exec(
		`INSERT INTO person (person_id, fullname, num_transaction, company_name) 
		VALUES ($1, $2, $3, $4)`,
		p.PersonId, p.FullName, p.NumTransaction, p.CompanyName)
	return err
}

func updatePerson(db *sql.DB, p PersonTarget) error {
	_, err := db.Exec(
		"UPDATE person SET fullname=$1, num_transaction=$2, company_name=$3 WHERE person_id=$4",
		p.FullName, p.NumTransaction, p.CompanyName, p.PersonId)
	return err
}

func deletePerson(db *sql.DB, personId uint) error {
	_, err := db.Exec(
		"DELETE FROM public.person where person_id =$1",
		personId)
	return err
}

func main() {
	mysqlConn, err := sql.Open("mysql", "root:password@(localhost:3307)/cdc_reference")
	if err != nil {
		log.Fatalln(err)
	}
	defer mysqlConn.Close()

	connStr := "user=admin password=password dbname=cdc_target sslmode=disable"
	pgConn, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalln(err)
	}
	defer pgConn.Close()

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "docker_mysql_8",
		"auto.offset.reset": "smallest"})
	if err != nil {
		log.Fatalf(err.Error())
	}

	log.Println("Ready")

	topics := []string{"docker_mysql_8.cdc_source.person"}
	err = consumer.SubscribeTopics(topics, nil)

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	run := true
	for run {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			ev, err := consumer.ReadMessage(100 * time.Millisecond)
			if err != nil {
				continue
			}
			eventValueRaw := ev.Value
			// eventKeyRaw := ev.Key
			// log.Println(string(eventKeyRaw))
			log.Println(string(eventValueRaw))
			if len(eventValueRaw) != 0 {
				var eventValue DebeziumValueEvent
				err = json.Unmarshal(eventValueRaw, &eventValue)
				if err != nil {
					log.Println(eventValueRaw)
					log.Fatalln(err)
				}

				switch eventValue.Payload.Operation {
				case "c":
					log.Println("Insert event")
					companyName, _ := getCompany(mysqlConn, eventValue.Payload.After.CompanyId)
					p := &PersonTarget{
						PersonId:       eventValue.Payload.After.PersonId,
						FullName:       eventValue.Payload.After.FullName,
						NumTransaction: eventValue.Payload.After.NumTransaction,
					}
					if companyName != nil {
						p.CompanyName = *companyName
					}

					err = insertPerson(pgConn, *p)
					if err != nil {
						log.Fatalln(err)
					}
				case "u":
					log.Println("Update event")
					companyName, _ := getCompany(mysqlConn, eventValue.Payload.After.CompanyId)
					p := &PersonTarget{
						PersonId:       eventValue.Payload.After.PersonId,
						FullName:       eventValue.Payload.After.FullName,
						NumTransaction: eventValue.Payload.After.NumTransaction,
					}
					if companyName != nil {
						p.CompanyName = *companyName
					}
					err = updatePerson(pgConn, *p)
					if err != nil {
						log.Fatalln(err)
					}
				case "d":
					log.Println("Delete event")
					deletePerson(pgConn, eventValue.Payload.Before.PersonId)
				default:
					log.Println("Unknown event")
					log.Println(string(eventValueRaw))
				}
			} else {
				// Tombstone events
			}
		}
	}

	consumer.Close()
}
