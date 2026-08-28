package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Nischaldh/Gator/internal/database"
	"github.com/google/uuid"
)

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %s <name> <url>", cmd.Name)
	}
	name := cmd.Args[0]
	url := cmd.Args[1]

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't create feed: %w", err)
	}
	follow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't follow feed: %w", err)
	}

	fmt.Println("Feed created successfully:")
	printFeed(feed)
	fmt.Println()
	fmt.Println("=====================================")
	fmt.Printf("%s is now following %s\n", follow.UserName, follow.FeedName)
	return nil

}

func printFeed(feed database.Feed) {
	fmt.Println("=====================================")
	fmt.Printf("* ID:            %s\n", feed.ID)
	fmt.Printf("* Created:       %v\n", feed.CreatedAt)
	fmt.Printf("* Updated:       %v\n", feed.UpdatedAt)
	fmt.Printf("* Name:          %s\n", feed.Name)
	fmt.Printf("* URL:           %s\n", feed.Url)
	fmt.Printf("* UserID:        %s\n", feed.UserID)
	fmt.Printf("* LastFetchedAt: %v\n", feed.LastFetchedAt.Time)
	fmt.Println("=====================================")
}


func handlerGetFeeds(s *state, cmd command) error{
	if len(cmd.Args)!=0{
		return fmt.Errorf("Too many arugments\n Usage %s", cmd.Name)
	}
	feeds, err:= s.db.GetFeeds(context.Background())
	if err!=nil{
		return fmt.Errorf("Error when fetching the feeds %v", err)
	}
	for _ , feed := range feeds{
		fmt.Println("=====================================")
		fmt.Printf("* Name:			%s\n",feed.Name)
		fmt.Printf("* URL:			%s\n",feed.Url)
		fmt.Printf("*User Name:		%s\n",feed.UserName)	
	}

	return nil
}


func handlerFollowFeed(s *state, cmd command, user database.User) error{
	if len(cmd.Args)!=1{
		return fmt.Errorf("Invlid formate/\n Usage %s <url>\n", cmd.Name)
	}
	feed, err:= s.db.GetFeed(context.Background(),cmd.Args[0])
	if err!=nil{
		return err
	}
	feeds, err:= s.db.CreateFeedFollow(context.Background(),database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err!=nil{
		return err
	}
	fmt.Println("===Feed Follow Created Successfully===")
	fmt.Println("=====================================")
	fmt.Printf("* Feed Name:			%s\n",feeds.FeedName)
	fmt.Printf("* User Name:			%s\n",feeds.UserName)
	return nil
}


func handlerFollowing(s *state, cmd command, user database.User) error{
	if len(cmd.Args)!=0{
		return fmt.Errorf("Too many arguments\n")
	}
	feeds, err:= s.db.GetFeedFollowsForUser(context.Background(),user.ID)
	if err!=nil{
		return err
	}
	fmt.Println("===FEEDS USER IS FOLLOWING===")
	for _, feed:=range feeds{
		fmt.Printf("%s\n", feed.FeedName)
	}
	return nil
}


func handlerDeleteFeedFollow(s *state, cmd command, user database.User) error{
	if len(cmd.Args)!=1{
		return fmt.Errorf("Invalid command.\nUsage %s <feed_url>", cmd.Name)
	}
	feed, err:= s.db.GetFeed(context.Background(), cmd.Args[0])
	if err!=nil{
		return err
	}
	if err:=s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	});err!=nil{
		return err
	}
	return nil
}