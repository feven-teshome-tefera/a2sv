package data

import (
	"context"
	"errors"
	"log"
	"time"

	"task_manager/db"
	"task_manager/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func getCollection() *mongo.Collection {
	client := db.Connect()
	return client.Database("tasksDB").Collection("tasks")
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
func GetTaskByID(id string) (models.Task, error) {
	collection := getCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Task{}, errors.New("invalid ID format")
	}

	var task models.Task
	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&task)
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
	res, err := collection.InsertOne(ctx, newTask)
	if err != nil {
		return err
	}

	newTask.ID = res.InsertedID.(primitive.ObjectID).Hex()
	log.Println("Inserted task with ID:", newTask.ID)
	return nil
}

func UpdateTask(id string, updated models.Task) error {
	collection := getCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid ID format")
	}

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
		update["dueDate"] = updated.DueDate
	}

	if len(update) == 0 {
		return errors.New("nothing to update")
	}

	res, err := collection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": update})
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

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid ID format")
	}

	res, err := collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("task not found")
	}

	return nil
}
