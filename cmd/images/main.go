package main

import (
	"fmt"

	"github.com/BajoJajoOrg/Inkscryption-backend/configs"
	"github.com/BajoJajoOrg/Inkscryption-backend/images/delivery"
	"github.com/BajoJajoOrg/Inkscryption-backend/images/usecase"
	"github.com/spf13/viper"
)

const (
	httpPath = "../../configs/image_http_config.yaml"
)

func main() {
	_, err := configs.LoadConfig(httpPath)
	if err != nil {
		fmt.Print("Ahtung")
	}

	psqInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		viper.Get("database.host"), viper.Get("database.port"), viper.Get("database.user"),
		viper.Get("database.password"), viper.Get("database.dbname"))

	fmt.Print((psqInfo))

	core, err := usecase.GetCore(psqInfo)
	if err != nil {
		fmt.Print("ZlukenSobaken")
		return
	}

	api := delivery.GetApi(core)

	errs := make(chan error, 2)

	go func() {
		errs <- api.ListenAndServe()
	}()

	err = <-errs
	if errs != nil {
		fmt.Printf("WEerr %s", err.Error())
	}
}
