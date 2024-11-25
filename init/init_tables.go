package initData

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/brianvoe/gofakeit/v7"
	"golibrary/internal/model"
	"golibrary/internal/service"
	"io"
	"net/http"
	"time"
)

func InitUsers(s *service.Service) error {
	users := make([]model.User, 0, 10)

	for i := 0; i < 10; i++ {
		gofakeit.Seed(time.Now().UnixNano())
		users = append(users, model.User{
			Name:       gofakeit.FirstName(),
			Surname:    gofakeit.LastName(),
			MiddleName: gofakeit.MiddleName(),
		})
	}

	for id, user := range users {
		age, err := GetUserAge(user.Name)
		if err != nil {
			return fmt.Errorf("init tables user [%d]: %w", id, err)
		}

		nation, err := GetUserNation(user.Name)
		if err != nil {
			return fmt.Errorf("init tables user [%d]: %w", id, err)
		}

		gender, err := GetUserGender(user.Name)
		if err != nil {
			return fmt.Errorf("init tables user [%d]: %w", id, err)
		}

		users[id].Age = age
		users[id].Gender = gender
		users[id].Nationality = nation

		_, err = s.User.CreateUser(context.Background(), &users[id])
		if err != nil {
			return fmt.Errorf("init tables user [%d]: %w", id, err)
		}
	}

	return nil
}

func Init(services *service.Service) error {
	err := InitUsers(services)
	if err != nil {
		return err
	}

	return nil
}

type Agify struct {
	Age int `json:"age"`
}

func GetUserAge(name string) (int, error) {
	resp, err := http.Get("https://api.agify.io?name=" + name)
	if err != nil {
		return 0, fmt.Errorf("Error to get user age: " + err.Error())
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("Error to get user age: " + err.Error())
	}

	var age Agify

	err = json.Unmarshal(data, &age)
	if err != nil {
		return 0, fmt.Errorf("Error to get user age: " + err.Error())
	}

	return age.Age, nil
}

type Nationalize struct {
	Country []Country `json:"country"`
}

type Country struct {
	Country string `json:"country_id"`
}

func GetUserNation(name string) (string, error) {
	resp, err := http.Get("https://api.nationalize.io/?name=" + name)
	if err != nil {
		return "", fmt.Errorf("Error to get user nation: " + err.Error())
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Error to get user nation: " + err.Error())
	}

	var nation Nationalize

	err = json.Unmarshal(data, &nation)
	if err != nil {
		return "", fmt.Errorf("Error to get user nation: " + err.Error())
	}

	if nation.Country == nil {
		return "", fmt.Errorf("Error to get user nation:", string(data))
	}

	return nation.Country[0].Country, nil
}

type Genderize struct {
	Gender string `json:"gender"`
}

func GetUserGender(name string) (string, error) {
	resp, err := http.Get("https://api.genderize.io?name=" + name)
	if err != nil {
		return "", fmt.Errorf("Error to get user gender: " + err.Error())
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Error to get user gender: " + err.Error())
	}

	var gender Genderize

	err = json.Unmarshal(data, &gender)
	if err != nil {
		return "", fmt.Errorf("Error to get user gender: " + err.Error())
	}

	return gender.Gender, nil
}
