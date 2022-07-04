package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"github.com/jamescostian/signal-to-sms/cmd"
)

var version, commit, date = "unknown", "unknown", "an unknown point in time"

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	err := cmd.Execute(ctx, fmt.Sprintf("%s (built from commit %s made in %s)", version, commit, date))
	// Only log errors that aren't directly caused by the user.
	// If a user hits Ctrl+C, they expect the program to stop, so no need to tell them that they caused an error.
	if err == nil || errors.Is(err, context.Canceled) {
		return
	}
	log.Fatalln("ERROR:", err)
}
