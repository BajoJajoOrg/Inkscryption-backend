package repo

import (
	"context"
	"database/sql"
	"fmt"
	"mime/multipart"
	"os"
	"time"

	structures "github.com/BajoJajoOrg/Inkscryption-backend/images"
	_ "github.com/jackc/pgx/stdlib"
	_ "github.com/lib/pq"

	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
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

func (storage *ImageStorage) Get(ctx context.Context, userID int64) ([]structures.Canvas, error) {
	//var images []image_struct.Image

	var canvases []structures.Canvas

	query := "SELECT " + canvasFields + " FROM canvas"

	stmt, err := storage.dbReader.Prepare(query)
	if err != nil {
		return []structures.Canvas{}, err
	}
	rows, err := stmt.Query()
	if err != nil {
		return []structures.Canvas{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var canvas structures.Canvas

		err = rows.Scan(&canvas.Name, &canvas.Url, &canvas.Update)
		if err != nil {
			return []structures.Canvas{}, err
		}

		canvases = append(canvases, canvas)
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		print("Error loading default config: %v", err)
		os.Exit(0)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(vkCloudHotboxEndpoint)
		o.Region = defaultRegion
	})

	presigner := s3.NewPresignClient(client)
	bucketName := "bajojajo"
	lifeTimeSeconds := int64(60)

	var req *v4.PresignedHTTPRequest

	var newCanvases []structures.Canvas

	for _, canvas := range canvases {
		objectKey := canvas.Url
		println("THIS IS OBJECT KEY", objectKey)
		req, err = presigner.PresignGetObject(context.TODO(), &s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(objectKey),
		}, func(opts *s3.PresignOptions) {
			opts.Expires = time.Duration(lifeTimeSeconds * int64(time.Second))
		})
		if err != nil {
			println(err.Error())
			return []structures.Canvas{}, err
		}
		newCanvas := structures.Canvas{
			Name:   canvas.Name,
			Url:    req.URL,
			Update: canvas.Update,
		}

		newCanvases = append(newCanvases, newCanvas)
	}

	return newCanvases, nil
}

// func (storage *ImageStorage) Add(ctx context.Context, img multipart.File) error {
// 	//logger := ctx.Value(Logg).(Log)
// 	query := "INSERT INTO canvas (canvas_name, url, update_time) VALUES ($1, $2, $3) ON CONFLICT (person_id, cell_number) DO UPDATE SET image_url = EXCLUDED.image_url;"

// 	logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("hehe ", image.UserId, image.CellNumber, image.Url)
// 	stmt, err := storage.dbReader.Prepare(query) // using prepared statement
// 	if err != nil {
// 		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("can't query: ", err.Error())
// 		return fmt.Errorf("Add img %w", err)
// 	}

// 	_, err = stmt.Exec(image.UserId, image.Url, image.CellNumber)
// 	//_, err := storage.dbReader.Exec(query, image.UserId, image.Url, image.CellNumber)
// 	if err != nil {
// 		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("can't query: ", err.Error())
// 		return fmt.Errorf("Add img %w", err)
// 	}

// 	sess, err := session.NewSession(&awsUpload.Config{
// 		Region: aws.String("ru-msk"),
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	svc := serviceUpload.New(sess, awsUpload.NewConfig().WithEndpoint(vkCloudHotboxEndpoint).WithRegion(defaultRegion))
// 	bucket := "los_ping"

// 	params := &serviceUpload.PutObjectInput{
// 		Bucket: aws.String(bucket),
// 		Key:    aws.String(image.FileName),
// 		Body:   img,
// 		ACL:    aws.String("public-read"),
// 	}

// 	_, err = svc.PutObject(params)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

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
	//_, err := storage.dbReader.Exec(query, image.UserId, image.Url, image.CellNumber)
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
		return err
	}
	return nil
}
