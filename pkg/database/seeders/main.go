package seeders

import (
	"context"
	"fmt"
	"go-rest/pkg/database/seeders/workers"
	"go.mongodb.org/mongo-driver/mongo"
	"reflect"
	"runtime"
	"time"
)

type SeederFunc func(ctx context.Context, client *mongo.Client, dbName string) error

func SeedDatabase(ctx context.Context, client *mongo.Client, seederType string, dbName string) error {
	var seeders []SeederFunc

	switch seederType {
	case "full":
		seeders = []SeederFunc{
			workers.SeedAdminUser,
			workers.SeedNotes,
		}
	case "basic":
		seeders = []SeederFunc{
			workers.SeedAdminUser,
		}
	default:
		return fmt.Errorf("unknown seeder type: %s", seederType)
	}

	for _, seeder := range seeders {
		seederName := getFunctionName(seeder)

		err := seeder(ctx, client, dbName)
		if err != nil {
			return fmt.Errorf("seeder %s failed: %w", seederName, err)
		}

		fmt.Printf("%s OK %s\n", time.Now().Format("2006/01/02 15:04:05"), seederName)
	}

	fmt.Println("Database seeding completed successfully.")
	return nil
}

func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
