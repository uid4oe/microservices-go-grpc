package advicedb

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
)

type Advice struct {
	UserId    string
	Advice    string
	CreatedAt time.Time
}

var _ = loadLocalEnv()
var (
	db       = GetEnv("POSTGRES_DB")
	username = GetEnv("POSTGRES_USER")
	password = GetEnv("POSTGRES_PASSWORD")
	host     = GetEnv("POSTGRES_HOST")
)

func NewClient(ctx context.Context) (*pgxpool.Pool, error) {
	url := "postgres://" + username + ":" + password + "@" + host + "/" + db
	client, err := pgxpool.Connect(ctx, url)
	if err != nil {
		return nil, errors.New("cannot connect to postgres instance")
	}
	return client, nil
}

func CreateOne(client *pgxpool.Pool, ctx context.Context, advice *Advice) error {
	_, err := client.Exec(ctx, "insert into advices(user_id,advice,created_at) values($1,$2,CURRENT_TIMESTAMP)", advice.UserId, advice.Advice)
	return err
}

func UpdateOne(client *pgxpool.Pool, ctx context.Context, advice *Advice) error {
	_, err := client.Exec(ctx, "update advices set advice=$1, created_at=CURRENT_TIMESTAMP where user_id=$2", advice.Advice, advice.UserId)
	return err
}

func FindOne(client *pgxpool.Pool, ctx context.Context, id string) (*Advice, error) {
	advice := Advice{UserId: id}
	err := client.QueryRow(ctx, "select advice,created_at from advices where user_id=$1", id).Scan(&advice.Advice, &advice.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &advice, nil
}

func loadLocalEnv() interface{} {
	if _, runningInContainer := os.LookupEnv("ADVICE_GRPC_SERVICE"); !runningInContainer {
		err := godotenv.Load("../.env.local")
		if err != nil {
			log.Fatal(err)
		}
	}
	return nil
}

func GetEnv(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		log.Fatal("Environment variable not found: ", key)
	}
	return value
}
