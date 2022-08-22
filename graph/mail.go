package graph

import (
	// graph "MZ221-TPA-Web-Back/graph"
	"MZ221-TPA-Web-Back/graph/model"
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/trycourier/courier-go/v2"
)

func SendActivationLink(r *Resolver, receiver *model.User) {
	var activationLink *model.Activation

	if err := r.DB.First(&activationLink, "user_id = ?", receiver.ID).Error; err != nil {
		activationLink = &model.Activation{
			ID:     uuid.NewString(),
			UserId: receiver.ID,
		}
		r.DB.Create(activationLink)
	}
	fmt.Println("2")

	// activationUrl := os.Getenv("FRONTEND_LINK") + "/" + activationLink.ID
	activationUrl := "http://localhost:5173" + "/activation/" + activationLink.ID
	fmt.Println(activationLink)

	// fmt.Println("courier token :" + os.Getenv("COURIER_AUTH_TOKEN"))
	// client := courier.CreateClient(os.Getenv("COURIER_AUTH_TOKEN"), nil)
	client := courier.CreateClient("pk_prod_XJKJ5HPAVPMQ0DMA822ZVE18XQ3T", nil)

	requestID, err := client.SendMessage(
		context.Background(),
		courier.SendMessageRequestBody{
			Message: map[string]interface{}{
				"to": map[string]string{
					"email": receiver.Email,
				},
				"template": "Z25DPSG3654FDHNVCFMHW8RM3Y3T",
				"data": map[string]string{
					"activationUrl": activationUrl,
				},
			},
		},
	)

	if err != nil {
		log.Fatalln(err)
	}
	log.Println(requestID)

}

func SendResetPasswordLink(r *Resolver, receiver *model.User, reset *model.Reset) {

	// resetUrl := os.Getenv("FRONTEND_LINK") + "/" + reset.ID
	resetUrl := "http://localhost:5173" + "/reset/" + reset.ID

	// client := courier.CreateClient(os.Getenv("COURIER_AUTH_TOKEN"), nil)
	client := courier.CreateClient("pk_prod_XJKJ5HPAVPMQ0DMA822ZVE18XQ3T", nil)

	requestID, err := client.SendMessage(
		context.Background(),
		courier.SendMessageRequestBody{
			Message: map[string]interface{}{
				"to": map[string]string{
					"email": receiver.Email,
				},
				"template": "2WGYHFDR4JMT1SNDWMQ3TDD6D48W",
				"data": map[string]string{
					"resetUrl": resetUrl,
				},
			},
		},
	)

	if err != nil {
		log.Fatalln(err)
	}
	log.Println(requestID)

}
