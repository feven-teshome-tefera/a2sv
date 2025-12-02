package data

import (
	"context"
	"errors"
	"fmt"
	"log"
	"task_manager/db"
	"task_manager/models"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	jwt "github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func getCollection() *mongo.Collection {
	client := db.Connect()
	return client.Database("tasksDB").Collection("tasks")
}
func getUserCollection() *mongo.Collection {
	client := db.Connect()
	return client.Database("tasksDB").Collection("users")
}

var JwtSecret = []byte("your_jwt_secret")

func LoginUser(c *gin.Context) (string, error) {
	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		return "", errors.New("invalid request")
	}

	collection := getUserCollection()

	var existingUser models.User
	err := collection.FindOne(context.Background(), bson.M{"email": input.Email}).Decode(&existingUser)
	if err != nil {
		return "", errors.New("invalid email or password")
	}
	if bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(input.Password)) != nil {
		return "", errors.New("invalid email or password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": existingUser.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})

	jwtToken, err := token.SignedString(JwtSecret)
	if err != nil {
		return "", errors.New("could not create token")
	}

	return jwtToken, nil
}
func GetAllTasks() ([]models.Task, error) {
	collection := getCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tasks []models.Task
	for cursor.Next(ctx) {
		var t models.Task
		if err := cursor.Decode(&t); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	return tasks, nil
}
func GetAlluser() ([]models.User, error) {
	collection := getUserCollection()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var user []models.User
	for cursor.Next(ctx) {
		var f models.User
		if err := cursor.Decode(&f); err != nil {
			return nil, err
		}
		user = append(user, f)
	}

	return user, nil
}
func GetTaskByID(id string) (models.Task, error) {
	collection := getCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var task models.Task
	err := collection.FindOne(ctx, bson.M{"id": id}).Decode(&task)
	if err != nil {
		return models.Task{}, errors.New("task not found")
	}

	return task, nil
}

func AddTask(newTask models.Task) error {
	if newTask.DueDate.IsZero() {
		newTask.DueDate = time.Now()
	}

	collection := getCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Generate custom ID
	count, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return err
	}
	newTask.ID = fmt.Sprintf("%d", count+1)

	_, err = collection.InsertOne(ctx, newTask)
	if err != nil {
		return err
	}

	log.Println("Inserted task with ID:", newTask.ID)
	return nil
}

func AddUser(newUser models.User) error {
	collection := getUserCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := collection.InsertOne(ctx, newUser)
	if err != nil {
		return err
	}

	log.Println("Inserted user with email:", newUser.Email)

	return nil
}

func UpdateTask(id string, updated models.Task) error {
	collection := getCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{}
	if updated.Title != "" {
		update["title"] = updated.Title
	}
	if updated.Description != "" {
		update["description"] = updated.Description
	}
	if updated.Status != "" {
		update["status"] = updated.Status
	}
	if !updated.DueDate.IsZero() {
		update["due_date"] = updated.DueDate
	}

	if len(update) == 0 {
		return errors.New("nothing to update")
	}

	res, err := collection.UpdateOne(ctx, bson.M{"id": id}, bson.M{"$set": update})
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("task not found")
	}

	return nil
}

func DeleteTask(id string) error {
	collection := getCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := collection.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("task not found")
	}

	return nil
}
func Deleteusers(email string) error {
	collection := getUserCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := collection.DeleteOne(ctx, bson.M{"email": email})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("user not found")
	}

	return nil
}
