package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/andyfusniak/monolith/internal/store/sqlite3"
	"github.com/andyfusniak/monolith/service"
)

func main() {
	// make a store
	rw, err := sqlite3.OpenDB("/tmp/monolith.db")
	if err != nil {
		log.Fatal(err)
	}
	defer rw.Close()
	rw.SetMaxOpenConns(1)
	rw.SetMaxIdleConns(1)
	rw.SetConnMaxIdleTime(5 * time.Minute)
	mystore := sqlite3.NewStore(rw, rw)

	// make the service
	svc := service.New(service.WithRepository(mystore))

	ctx := context.Background()

	// user, err := svc.GetUser(ctx, "pBFkQiKQ3pUXiekGsDGuUZ")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("%#v\n", user)

	user, err := svc.VerifyUserPassword(ctx, "sawasdiworn.k@gmail.com", "testtesttest")
	if err != nil {
		fmt.Printf("VerifyUserPassword: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("%#v\n", user)

	// row, err := mystore.GetUser(context.Background(), "pBFkQiKQ3pUXiekGsDGuUZ")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Printf("%#v\n", row)
}
