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

type User struct {
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

type BaseRepository struct {
	collection *mongo.Collection
}

type userRepository struct {
	BaseRepository
}

type friendRepository struct {
	BaseRepository
}

type messageRepository struct {
	BaseRepository
}

type UserRepository interface {
	Insert(ctx context.Context, user User) error
	FindByID(ctx context.Context, id primitive.ObjectID) (User, error)
}

type MessageRepository interface {
	Insert(ctx context.Context, message Messages) error
}

type FriendsRepository interface {
	Insert(ctx context.Context, friend Friends) error
	FriendsByTimeRange(ctx context.Context, user User, startTime, endTime time.Time) ([]Friends, error)
	FriendsByUser(ctx context.Context, user User, userRepo *userRepository) ([]string, error)
	CountFriends(ctx context.Context, user User) (int64, error)
}

func NewUserRepository(client *mongo.Client) *userRepository {
	return &userRepository{
		BaseRepository: BaseRepository{
			collection: client.Database("cheeper").Collection("users"),
		},
	}
}

func NewMessageRepository(client *mongo.Client) *messageRepository {
	return &messageRepository{
		BaseRepository: BaseRepository{
			collection: client.Database("cheeper").Collection("messages"),
		},
	}
}

func NewFriendRepository(client *mongo.Client) *friendRepository {
	return &friendRepository{
		BaseRepository: BaseRepository{
			collection: client.Database("cheeper").Collection("friends"),
		},
	}
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

func (r *userRepository) Insert(ctx context.Context, user User) error {
	_, err := r.collection.InsertOne(ctx, user)
	return err
}

func (r *userRepository) FindByID(ctx context.Context, id primitive.ObjectID) (User, error) {
	filter := bson.M{
		"_id": id,
	}

	var user User
	err := r.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (r *messageRepository) Insert(ctx context.Context, message Messages) error {
	_, err := r.collection.InsertOne(ctx, message)
	return err
}

func (r *friendRepository) Insert(ctx context.Context, friend Friends) error {
	_, err := r.collection.InsertOne(ctx, friend)
	return err
}

func (r *friendRepository) FriendsByTimeRange(ctx context.Context, user User, startTime, endTime time.Time) ([]Friends, error) {
	filter := bson.M{
		"startDate": bson.M{
			"$gte": startTime,
			"$lte": endTime,
		},
		"user1": user.ID,
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var friendships []Friends
	if err = cursor.All(ctx, &friendships); err != nil {
		return nil, err
	}

	return friendships, nil
}

func (r *friendRepository) FriendsByUser(ctx context.Context, user User, userRepo *userRepository) ([]string, error) {
	opts := options.Find().SetSort(bson.M{"name": 1})
	filter := bson.M{
		"user1": user.ID,
	}

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var frndNames []string
	for cursor.Next(ctx) {
		var friend Friends
		if err = cursor.Decode(&friend); err != nil {
			return nil, err
		}
		usr, err := userRepo.FindByID(ctx, friend.User2)
		if err != nil {
			return nil, err
		}
		frndNames = append(frndNames, usr.Name)
	}

	return frndNames, nil

}

func (r *friendRepository) CountFriends(ctx context.Context, user User) (int64, error) {
	filter := bson.M{
		"user1": user.ID,
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func PrintFriends(friendships *[]Friends) {
	for _, frnd := range *friendships {
		fmt.Printf("User1: %v, user2: %v, startDate: %v\n", frnd.User1, frnd.User2, frnd.StartDate)
	}
}

func main() {
	//Подключение к базе данных
	client, err := connectDB()
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных")
	}

	userRepo := NewUserRepository(client)
	messageRepo := NewMessageRepository(client)
	friendRepo := NewFriendRepository(client)

	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	hexUser1, _ := primitive.ObjectIDFromHex("681c82fe89aff140d966aeb9")
	hexUser2, _ := primitive.ObjectIDFromHex("681c82fe89aff140d966aebd")

	//Добавление нового пользователя, сообщения, друга
	newUser := User{
		ID:    primitive.NewObjectID(),
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

	err = userRepo.Insert(ctx, newUser)
	if err != nil {
		log.Println(err)
	}

	err = messageRepo.Insert(ctx, newMessage)
	if err != nil {
		log.Println(err)
	}

	err = friendRepo.Insert(ctx, newFriend)
	if err != nil {
		log.Println(err)
	}

	user := User{
		ID:    hexUser1,
		Name:  "name1",
		Login: "login1",
	}

	//Получение друзей у заданного пользователя за последние 2 часа
	endTime := time.Now()
	startTime := endTime.Add(-24 * time.Hour)
	friendships, err := friendRepo.FriendsByTimeRange(ctx, user, startTime, endTime)
	if err != nil {
		log.Printf("Ошибка при получении друзей: %v", err)
	}
	PrintFriends(&friendships)

	//Получение упорядоченного списка имен друзей заданного пользователя
	frndNames, err := friendRepo.FriendsByUser(ctx, user, userRepo)
	if err != nil {
		log.Print(err)
	}

	for _, name := range frndNames {
		fmt.Println(name)
	}

	//Вычисление количества друзей для заданного пользователя
	count, err := friendRepo.CountFriends(ctx, user)
	if err != nil {
		log.Print(err)
	}

	fmt.Printf("У пользователя %v %d друзей\n", user.Login, count)
}
