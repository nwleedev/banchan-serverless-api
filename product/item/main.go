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
type Request events.APIGatewayProxyRequest
type Response events.APIGatewayProxyResponse
type ResponseData struct {
	*models.Product
	Poster models.PosterSlice `json:"posters"`
	Tag    models.TagSlice    `json:"tags"`
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

	pathID := request.PathParameters["id"]
	id, err := strconv.Atoi(pathID)

	if err != nil || id < 1 {
		log.Println(err)
		return Response{StatusCode: 501, Body: `{"message": "Parameter Error"}`}, err
	}

	data, err := models.Products(
		qm.Load(models.ProductRels.Tags, qm.Where("product_id = ?", id)),
		qm.Load(models.ProductRels.Posters, qm.Where("product_id = ?", id)),
		qm.Where("id = ?", id),
	).One(ctx, db)

	product := &ResponseData{data, data.R.Posters, data.R.Tags}

	if err != nil {
		log.Println(err)
		return Response{StatusCode: 501, Body: `{"message": "DB Connection Error"}`}, err
	}

	body, err := json.Marshal(map[string]interface{}{
		"message": "OK",
		"product": product,
	})

	if err != nil {
		log.Println(err)
		return Response{StatusCode: 501}, err
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
