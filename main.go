package main

import (
	"context"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"iter"
	"net/http"
	"time"
)

const Window = 10000

type Entity struct {
	ID    int64
	Code  string
	Value int
}

func getEntities(db *gorm.DB) iter.Seq[Entity] {
	return func(yield func(Entity) bool) {
		offset := 0
		count := 1

		for count != 0 {
			entities, err := gorm.G[Entity](db).Limit(Window).Offset(offset).Order("id asc").Find(context.Background())
			if err != nil {
				return
			}

			for _, v := range entities {
				if !yield(v) {
					return
				}
			}

			count = len(entities)
			offset += Window
			time.Sleep(time.Second * 1)
		}
	}
}

func main() {
	dsn := "host=localhost user=playground_user password=playground_password dbname=playground_db port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/entities", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		for v := range getEntities(db) {
			row := fmt.Sprintf("%d,%s,%d\n", v.ID, v.Code, v.Value)
			fmt.Printf("write row -> %s", row)
			w.Write([]byte(row))
		}
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
