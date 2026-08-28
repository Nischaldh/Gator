package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Nischaldh/Gator/internal/database"
	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %v <name>", cmd.Name)
	}
	_, err := s.db.GetUser(context.Background(), cmd.Args[0])
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return fmt.Errorf("user does not exist")
        }

        return err
    }

    err = s.cfg.SetUser(cmd.Args[0])
    if err != nil {
        return err
    }

    fmt.Printf("User has been set to %s\n", cmd.Args[0])
    return nil
}


func handlerRegister(s *state, cmd command) error{
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Username is required")
	}
	name := cmd.Args[0]
	ctx := context.Background()
	_, err:= s.db.GetUser(ctx, name)
	if err == nil{
		return fmt.Errorf("User already exists")
	}
	if !errors.Is(err, sql.ErrNoRows){
		return err
	}
	user, err := s.db.CreateUser(ctx, database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: name,
	})
	if err != nil{
		return err
	}
	err = s.cfg.SetUser(name)
	if err!=nil{
		return err
	}
	fmt.Printf("User %s was created\n", name)
	fmt.Printf("%v\n", user)
	return nil
}


func handlerReset(s *state, cmd command) error{
	err:=s.db.Reset(context.Background())
	if err !=nil{
		return fmt.Errorf("Error when reseting %v", err)
	}
	fmt.Println("Successfully reset.")
	return nil
}

func handlerUsers(s *state, cmd command) error{
	users, err := s.db.GetUsers(context.Background())
	if err != nil{
		return err
	}
	for _, user := range users{
		if (user.Name == s.cfg.CurrentUserName){
			fmt.Printf("* %s (current)\n", user.Name)
		}else{
			fmt.Printf("* %s \n", user.Name)
		}
	} 
	return nil
}