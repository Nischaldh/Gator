package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Nischaldh/Gator/internal/database"
	"github.com/google/uuid"
)

func scrapeFeeds(s *state) error {
	next_feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}
	fmt.Printf("Fetching feed: %s\n", next_feed.Name)
	if err := s.db.MarkFeedFetched(context.Background(), next_feed.ID); err != nil {
		return fmt.Errorf("Error while fetching the feed: %w\n", err)
	}
	rssFeed, err := fetchFeed(context.Background(), next_feed.Url)
	if err != nil {
		return err
	}
	for _, item := range rssFeed.Channel.Item {
		publishedAt := sql.NullTime{}
		if t, err := time.Parse(time.RFC1123Z, item.PubDate); err == nil {
			publishedAt = sql.NullTime{
				Time:  t,
				Valid: true,
			}
		}
		_, err = s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			FeedID:    next_feed.ID,
			Title:     item.Title,
			Description: sql.NullString{
				String: item.Description,
				Valid:  true,
			},
			Url:         item.Link,
			PublishedAt: publishedAt,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			log.Printf("Couldn't create post: %v", err)
			continue
		}

	}
	log.Printf("Feed %s collected, %v posts found", next_feed.Name, len(rssFeed.Channel.Item))

	return nil

}
