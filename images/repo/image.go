package repo

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"mime/multipart"
	"time"

	structures "github.com/BajoJajoOrg/Inkscryption-backend/images"
	_ "github.com/jackc/pgx/stdlib"
	_ "github.com/lib/pq"

	"github.com/aws/aws-sdk-go/aws"
	awsUpload "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"

	//"github.com/aws/aws-sdk-go/aws/session"
	serviceUpload "github.com/aws/aws-sdk-go/service/s3"
)

const (
	vkCloudHotboxEndpoint = "https://hb.ru-msk.vkcs.cloud"
	defaultRegion         = "ru-msk"
)

type ImageStorage struct {
	dbReader *sql.DB
}

const (
	personImageFields = "person_id, image_url"
	canvasFields      = "canvas_name, url, update_time"
)

func NewImageStorage(dbReader *sql.DB) *ImageStorage {
	return &ImageStorage{
		dbReader: dbReader,
	}
}

func GetImageRepo(config string) (*ImageStorage, error) {
	db, err := sql.Open("postgres", config)
	if err != nil {
		println(err.Error())
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(90)
	db.SetMaxIdleConns(90)
	if err = db.Ping(); err != nil {
		println(err.Error())
	}

	postgreDb := ImageStorage{dbReader: db}

	go postgreDb.pingDb(50)
	return &postgreDb, nil
}

func (storage *ImageStorage) pingDb(timer uint32) {
	for {
		err := storage.dbReader.Ping()
		if err != nil {
			println(err.Error())
		}
		print("pong")

		time.Sleep(time.Duration(timer) * time.Second)
	}
}

func (storage *ImageStorage) Get(ctx context.Context, userID int64, dates []string) ([]structures.Canvas, error) {
	//var images []image_struct.Image

	var canvases []structures.Canvas

	query := "SELECT " + canvasFields + " FROM canvas"

	var args []interface{}

	fmt.Print("\nThisisdates\n", dates, "\n")

	if len(dates) != 0 {
		//query += " WHERE update_time BETWEEN " + dates[0] + " AND " + dates[1]
		query += " WHERE update_time BETWEEN $1 AND $2"
		args = append(args, dates[0], dates[1])
	}

	rows, err := storage.dbReader.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var canvas structures.Canvas

		err = rows.Scan(&canvas.Name, &canvas.Url, &canvas.Update)
		if err != nil {
			return nil, err
		}

		canvases = append(canvases, canvas)
	}

	return canvases, nil
}

func (storage *ImageStorage) Add(ctx context.Context, canvas structures.Canvas, img multipart.File) error {
	//logger := ctx.Value(Logg).(Log)
	query := "INSERT INTO canvas (canvas_name, url, update_time) VALUES ($1, $2, $3);"

	//logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("hehe ", image.UserId, image.CellNumber, image.Url)
	stmt, err := storage.dbReader.Prepare(query) // using prepared statement
	if err != nil {
		//logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("can't query: ", err.Error())
		return fmt.Errorf("Add img %w", err)
	}

	_, err = stmt.Exec(canvas.Name, canvas.Url, canvas.Update)
	if err != nil {
		//logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("can't query: ", err.Error())
		return fmt.Errorf("Add img %w", err)
	}

	sess, err := session.NewSession(&awsUpload.Config{
		Region: aws.String("ru-msk"),
	})
	if err != nil {
		return err
	}

	svc := serviceUpload.New(sess, awsUpload.NewConfig().WithEndpoint(vkCloudHotboxEndpoint).WithRegion(defaultRegion))
	bucket := "bajojajo"

	params := &serviceUpload.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(canvas.Name),
		Body:   img,
		ACL:    aws.String("public-read"),
	}

	_, err = svc.PutObject(params)
	if err != nil {
		log.Fatal("erorr!!", err)
		return err
	}
	log.Print("something is happening")
	return nil
}

func (storage *ImageStorage) Update(ctx context.Context, canvas structures.Canvas, img multipart.File) error {
	// //logger := ctx.Value(Logg).(Log)
	query := `UPDATE canvas
			SET update_time = $1
			WHERE canvas_name = $2`

	//logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("hehe ", image.UserId, image.CellNumber, image.Url)
	stmt, err := storage.dbReader.Prepare(query) // using prepared statement
	if err != nil {
		//logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("can't query: ", err.Error())
		return fmt.Errorf("Add img %w", err)
	}

	_, err = stmt.Exec(time.Now(), canvas.Name)
	if err != nil {
		//logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("can't query: ", err.Error())
		return fmt.Errorf("Add img %w", err)
	}

	sess, err := session.NewSession(&awsUpload.Config{
		Region: aws.String("ru-msk"),
	})
	if err != nil {
		return err
	}

	svc := serviceUpload.New(sess, awsUpload.NewConfig().WithEndpoint(vkCloudHotboxEndpoint).WithRegion(defaultRegion))
	bucket := "bajojajo"

	params := &serviceUpload.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(canvas.Name),
		Body:   img,
		ACL:    aws.String("public-read"),
	}

	_, err = svc.PutObject(params)
	if err != nil {
		log.Fatal("erorr!!", err)
		return err
	}
	log.Print("something is happening")
	return nil
}

func (storage *ImageStorage) Delete(ctx context.Context, canvas structures.Canvas) error {
	query := "DELETE FROM canvas WHERE canvas_name = $1"

	fmt.Print(canvas.Name)

	_, err := storage.dbReader.Exec(query, canvas.Name)
	if err != nil {
		//log.Fatalf("fatal %w", err)
		fmt.Printf("error in db", err)
		return fmt.Errorf("Delete img %w", err)
	}

	sess, err := session.NewSession(&awsUpload.Config{
		Region: aws.String("ru-msk"),
	})
	if err != nil {
		return err
	}

	svc := serviceUpload.New(sess, awsUpload.NewConfig().WithEndpoint(vkCloudHotboxEndpoint).WithRegion(defaultRegion))
	bucket := "bajojajo"
	key := "1/" + canvas.Name

	input := &serviceUpload.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(key),
	}
	result, err := svc.ListObjectsV2(input)
	if err != nil {
		log.Fatalf("Unable to list objects in directory %q, %v\n", key, err)
	}

	for _, obj := range result.Contents {
		if _, err := svc.DeleteObject(&serviceUpload.DeleteObjectInput{
			Bucket: aws.String(bucket),
			Key:    obj.Key,
		}); err != nil {
			log.Fatalf("Unable to delete object %q from bucket %q, %v\n", key, bucket, err)
		} else {
			log.Printf("Object %q deleted from bucket %q\n", key, bucket)
		}
	}
	return nil

}
