package main

import (
	//"flag"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/urfave/cli"
	"os"
	"strings"
)

func Database() *bolt.DB {
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		panic(err)
	}

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("TASK"))
		if err != nil {
			fmt.Errorf("error: %v", err)
		}

		return nil
	})

	return db
}

func Add(c *cli.Context) error {
	db := Database()
	args := strings.Join(c.Args(), " ")
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("TASK"))

		err := b.Put([]byte(args), []byte("TODO"))
		if err != nil {
			fmt.Errorf("error: %v", err)
		}
		return nil
	})

	return nil
}

func Do(c *cli.Context) error {
	db := Database()
	args := strings.Join(c.Args(), " ")
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("TASK"))
		v := b.Get([]byte(args))
		if string(v) == "TODO" {
			err := b.Put([]byte(args), []byte("DONE"))
			if err != nil {
				fmt.Errorf("error: %v", err)
			}
		}

		return nil
	})
	return nil
}

func List(c *cli.Context, s string) error {
	db := Database()
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("TASK"))
		i := 1
		b.ForEach(func(k, v []byte) error {
			if string(v) == s {
				fmt.Printf("%d. %s\n", i, string(k))
				i++
			}
			return nil
		})

		return nil
	})
	return nil
}

func Completed(c *cli.Context) error {
	return List(c, "DONE")
}

func Incompleted(c *cli.Context) error {
	return List(c, "TODO")
}

func main() {
	app := cli.NewApp()
	app.Name = "task"
	app.Usage = "task is a CLI for managing your TODOs."
	app.Commands = []cli.Command{
		{
			Name:   "add",
			Usage:  "Add a new task to your TODO list.",
			Action: Add,
		},
		{
			Name:   "do",
			Usage:  "Mark a task on your TODO list as complete.",
			Action: Do,
		},
		{
			Name:   "list",
			Usage:  "List all of your incomplete tasks.",
			Action: Incompleted,
		},
		{
			Name:   "completed",
			Usage:  "List all of your complete tasks.",
			Action: Completed,
		},
	}
	app.Run(os.Args)
}
