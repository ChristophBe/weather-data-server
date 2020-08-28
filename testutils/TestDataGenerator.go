package testing

import (
	"github.com/ChristophBe/weather-data-server/data/database"
	"github.com/ChristophBe/weather-data-server/data/models"
	"golang.org/x/tools/go/ssa/interp/testdata/src/fmt"
	"math/rand"
	"testing"
	"time"
)

func GetSavedMeasuringNode(t *testing.T) (models.MeasuringNode, models.User)  {

	node:= models.MeasuringNode{

		Name:       generateRandomString(),
		Lat:        rand.Float64(),
		Lng:        rand.Float64(),
		IsPublic:   false,
		IsOutdoors: false,
	}

	owner:= GetSavedUser(t)
	repo := database.GetMeasuringNodeRepository()
	node, err := repo.CreateMeasuringNode(node,owner.Id)
	if err != nil{
		t.Fatalf("can not create saved node; cause: %f",err)
	}
	return node, owner
}

func GetSavedUser(t *testing.T) models.User {
	user:= models.User{
		LastLogin:    time.Time{},
		CreationTime: time.Time{},
		Email:        generateRandomEmail(),
		Username:     generateRandomString(10),
		IsEnabled:    true,
		PasswordHash: []byte(generateRandomString(64)),
	}
	user, err  := database.GetUserRepository().SaveUser(user)
	if err != nil{
		t.Fatalf("can not create saved node; cause: %f",err)
	}
	return user
}

func generateRandomEmail() string {
	return fmt.Sprintf("%s@%s",generateRandomString(10),generateRandomString(10))
}

func generateRandomString(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}