package storage

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"

	utils "github.com/didikizi/RedSoft/packege"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	pool *pgxpool.Pool
	humanQuery
	imageQuery
}

type imageQuery struct {
	insertMail string
	selectMail string
	deleteMail string
}

type humanQuery struct {
	selectAllHuman         string
	selectHumanFromId      string
	insertHuman            string
	deleteHuman            string
	selectHumanFromSurname string
	updateHumanFromId      string
}

type Config interface {
	GetStoragePort() string
	GetStorageAddr() string
	GetStorageBaseName() string
	GetStorageUserLogin() string
	GetStorageUserPass() string
}

func New(ctx context.Context, config Config) (*Storage, error) {
	storageConnString := fmt.Sprintf("postgres://%s:%s@%s/%s",
		config.GetStorageUserLogin(),
		config.GetStorageUserPass(),
		net.JoinHostPort(config.GetStorageAddr(), config.GetStoragePort()),
		config.GetStorageBaseName(),
	)

	fmt.Println(storageConnString)

	database, err := pgxpool.New(ctx, storageConnString)
	if err != nil {
		slog.Info(utils.GetCallerInfo(), slog.String("err:", err.Error()))
		return &Storage{}, err
	}

	err = database.Ping(ctx)
	if err != nil {
		slog.Info(utils.GetCallerInfo(), slog.String("err:", err.Error()))
		return &Storage{}, err
	}

	err = migrate(ctx, database)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &Storage{pool: database,
		humanQuery: humanQuery{
			selectAllHuman:         "Select id, age, status, surname, name, fatherland, national, sex From human;",
			selectHumanFromId:      "Select id, age, status, surname, name, fatherland, national, sex From human Where id = @id;",
			selectHumanFromSurname: "Select id, age, status, surname, name, fatherland, national, sex From human Where surname = @surname;",
			insertHuman: `INSERT INTO human (age, status, surname, name, fatherland, national, sex)
				VALUES (@age, @status, @surname, @name, @fatherland, @national, @sex);`,
			updateHumanFromId: `UPDATE human SET age = @age, surname = @surname, name = @name, fatherland = @fatherland, national = @national, sex = @sex WHERE id = @id;`,
			deleteHuman:       "Delete From human Where id = @id;",
		},
		imageQuery: imageQuery{
			insertMail: `INSERT INTO mails (human_id, mail, description)
    			VALUES (@human_id, @mail, @description);`,
			selectMail: "Select id, human_id, mail, description From mails Where human_id = @human_id;",
			deleteMail: "Delete From mails Where id = @id;",
		}}, nil
}

func migrate(ctx context.Context, db *pgxpool.Pool) error {
	files, err := os.ReadDir("iternal/migrations")
	if err != nil {
		return fmt.Errorf("read migrations dir error: %s", err)
	}

	migrations := []string{}

	if len(files) < 1 {
		return fmt.Errorf("migrations not found")
	}

	for _, v := range files {
		content, errRead := os.ReadFile(fmt.Sprint("iternal/migrations/", v.Name()))
		if errRead != nil {
			return fmt.Errorf("failed to read migration file: %s, filename: %s", errRead, v.Name())
		}

		migrations = append(migrations, string(content))
	}

	for _, m := range migrations {
		_, err := db.Exec(ctx, m)
		if err != nil {
			return fmt.Errorf("%s fail query migration", err)
		}
	}

	return err
}
