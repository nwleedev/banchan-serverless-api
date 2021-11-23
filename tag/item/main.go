package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/projects/banchan-app/models"
	"github.com/projects/banchan-app/pkg"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse
type Request events.APIGatewayProxyRequest
type ResponseData struct {
	*models.Tag
	Products models.ProductSlice `json:"products"`
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, request Request) (Response, error) {
	var buf bytes.Buffer

	conn := new(pkg.Conn)
	dbCtx, db := conn.GetConnection()
	defer dbCtx.Done()

	if db == nil {
		return Response{StatusCode: 501, Body: `{"message": "DB Connection Error"}`}, errors.New("DB Connection Error")
	}

	defer db.Close()

	_id := request.PathParameters["id"]
	id, err := strconv.Atoi(_id)
	if err != nil || id < 1 {
		log.Println(err)
		return Response{StatusCode: 501, Body: `{"message": "Parameter ID Error"}`}, err
	}

	_productID := request.PathParameters["product_id"]
	productID, err := strconv.Atoi(_productID)
	if err != nil || id < 1 {
		log.Println(err)
		return Response{StatusCode: 501, Body: `{"message": "Parameter Product ID Error"}`}, err
	}

	data, err := models.Tags(
		qm.Where("id = ?", id),
		qm.Load(
			models.TagRels.Products,
			qm.OrderBy("product_id desc"),
			qm.Where("tag_id = ?", id),
			qm.Limit(20),
			qm.Offset((productID-1)*20),
		),
	).One(dbCtx, db)

	tag := &ResponseData{
		data,
		data.R.Products,
	}

	if err != nil {
		log.Println(err)
		return Response{StatusCode: 501, Body: `{"message": "DB Connection Error"}`}, err
	}

	body, err := json.Marshal(map[string]interface{}{
		"message": "OK",
		"tag":     tag,
	})
	if err != nil {
		log.Println(err)
		return Response{StatusCode: 404}, err
	}
	json.HTMLEscape(&buf, body)

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type":                     "application/json",
			"Access-Control-Allow-Origin":      "*",
			"Access-Control-Allow-Credentials": "true",
			"Access-Control-Allow-Methods":     "GET",
		},
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
