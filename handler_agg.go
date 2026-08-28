package main

import (
	"fmt"
	"time"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args)!=1{
		return fmt.Errorf("usage: agg <time_between_reqs>")
	}
	time_between_request, err:= time.ParseDuration(cmd.Args[0])
	if err!=nil{
		return fmt.Errorf("invalid duration: %w", err)
	}
		fmt.Printf(
		"Collecting feeds every %v\n",
		time_between_request,
	)
	ticker := time.NewTicker(time_between_request)
	defer ticker.Stop()
	for ; ; <-ticker.C {
		err := scrapeFeeds(s)
		if err != nil {
			fmt.Printf("Error scraping feeds: %v\n", err)
		}
	}

}


