package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/jackc/pgx"
)

func main() {
	flag.Parse()
	erf := func(err error) {
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
	cstr, err := pgx.ParseURI(os.Getenv("DATABASE_URL"))
	erf(err)
	conn, err := pgx.Connect(cstr)
	erf(err)

	os.Args[0] = strings.TrimPrefix(os.Args[0], "./")
	os.Args[0] = strings.TrimPrefix(os.Args[0], "sql")

	q, err := conn.Query(strings.Join(os.Args, " "))
	erf(err)
	for _, d := range q.FieldDescriptions() {
		fmt.Print(string(d.Name), "\t")
	}
	fmt.Println()

	for q.Next() {
		vs, err := q.Values()
		erf(err)
		for _, v := range vs {
			fmt.Print(v, "\t")
		}
		fmt.Println()
	}
	if err := q.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	conn.Close()
}
