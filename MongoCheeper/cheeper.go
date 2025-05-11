package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Insertion interface {
	Insert(client *mongo.Client)
}

type Users struct {
	ID    primitive.ObjectID `bson:"_id"`
	Name  string             `bson:"name"`
	Login string             `bson:"login"`
}

type Friends struct {
	User1     primitive.ObjectID `bson:"user1"`
	User2     primitive.ObjectID `bson:"user2"`
	StartDate time.Time          `bson:"startDate"`
}

type Messages struct {
	Sender    primitive.ObjectID `bson:"sender"`
	Recipient primitive.ObjectID `bson:"recipient"`
	Text      string             `bson:"text"`
}

func connectDB() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func Insert(client *mongo.Client, insertion Insertion) {
	insertion.Insert(client)
}

func (u Users) Insert(client *mongo.Client) {
	collection := client.Database("cheeper").Collection("users")
	_, err := collection.InsertOne(context.TODO(), u)
	if err != nil {
		log.Printf("Ошибка при вставке пользователя: %v", err)
	}
}

func (m Messages) Insert(client *mongo.Client) {
	collection := client.Database("cheeper").Collection("messages")
	_, err := collection.InsertOne(context.TODO(), m)
	if err != nil {
		log.Printf("Ошибка при вставке сообщения: %v", err)
	}
}

func (f Friends) Insert(client *mongo.Client) {
	collection := client.Database("cheeper").Collection("friends")
	_, err := collection.InsertOne(context.TODO(), f)
	if err != nil {
		log.Printf("Ошибка при добавлении новой дружбы: %v", err)
	}
}

func getFriendsByTimeRange(client *mongo.Client, user Users, startTime, endTime time.Time) ([]Friends, error) {
	collection := client.Database("cheeper").Collection("friends")

	filter := bson.M{
		"startDate": bson.M{
			"$gte": startTime,
			"$lte": endTime,
		},
		"user1": user.ID,
	}

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var friendships []Friends
	if err = cursor.All(context.TODO(), &friendships); err != nil {
		return nil, err
	}

	return friendships, nil
}

func findUser(client *mongo.Client, friend Friends) (Users, error) {
	usersCollection := client.Database("cheeper").Collection("users")

	filter := bson.M{
		"_id": friend.User2,
	}

	var user Users
	err := usersCollection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return Users{}, err
	}

	return user, nil

}

func getFriendsByUser(client *mongo.Client, user Users) ([]string, error) {
	friendsCollection := client.Database("cheeper").Collection("friends")

	opts := options.Find().SetSort(bson.M{"age": 1})
	filter := bson.M{
		"user1": user.ID,
	}

	cursor, err := friendsCollection.Find(context.TODO(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var frndNames []string
	for cursor.Next(context.TODO()) {
		var friend Friends
		if err = cursor.Decode(&friend); err != nil {
			return nil, err
		}
		usr, err := findUser(client, friend)
		if err != nil {
			return nil, err
		}
		frndNames = append(frndNames, usr.Name)
	}

	return frndNames, nil

}

func countFriends(client *mongo.Client, user Users) (int, error) {
	friends, err := getFriendsByUser(client, user)
	if err != nil {
		return 0, err
	}
	return len(friends), nil
}

func PrintFrnd(friendships *[]Friends) {
	for _, frnd := range *friendships {
		fmt.Printf("User1: %v, user2: %v, startDate: %v\n", frnd.User1, frnd.User2, frnd.StartDate)
	}
}

func main() {
	//Подключение к базе данных
	client, _ := connectDB()
	if client == nil {
		log.Fatal("Ошибка подключения к базе данных")
	}

	hexUser1, _ := primitive.ObjectIDFromHex("681c82fe89aff140d966aeb9")
	hexUser2, _ := primitive.ObjectIDFromHex("681c82fe89aff140d966aebd")

	//Добавление нового пользователя, сообщения, друга
	newUser := Users{
		Name:  "NewName11",
		Login: "NewLogin11",
	}

	newMessage := Messages{
		Sender:    hexUser1,
		Recipient: hexUser2,
		Text:      "New message",
	}

	newFriend := Friends{
		User1:     hexUser1,
		User2:     hexUser2,
		StartDate: time.Now(),
	}

	Insert(client, newUser)
	Insert(client, newMessage)
	Insert(client, newFriend)

	user := Users{
		ID:    hexUser1,
		Name:  "name1",
		Login: "login1",
	}

	//Получение друзей у заданного пользователя за последние 2 часа
	endTime := time.Now()
	startTime := endTime.Add(-24 * time.Hour)
	friendships, err := getFriendsByTimeRange(client, user, startTime, endTime)
	if err != nil {
		log.Printf("Ошибка при получении друзей: %v", err)
	}
	PrintFrnd(&friendships)

	//Получение упорядоченного списка имен друзей заданного пользователя
	frndNames, err := getFriendsByUser(client, user)
	if err != nil {
		log.Print(err)
	}

	for _, name := range frndNames {
		fmt.Println(name)
	}

	//Вычисление количества друзей для заданного пользователя
	count, err := countFriends(client, user)
	if err != nil {
		log.Print(err)
	}

	fmt.Printf("У пользователя %v %d друзей\n", user.Login, count)
}
