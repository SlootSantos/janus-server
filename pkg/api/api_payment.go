// +build !enterprise

package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/SlootSantos/janus-server/pkg/api/auth"
	"github.com/SlootSantos/janus-server/pkg/storage"
	"github.com/stripe/stripe-go/v71"
	"github.com/stripe/stripe-go/v71/checkout/session"
)

var _ = setupPayment()

func setupPayment() bool {
	log.Println("PAYMENT IS ACTIVE!")
	os.Setenv("IS_ENTERPRISE", "false")

	http.HandleFunc("/payment/success/", auth.WithCredentials(func(w http.ResponseWriter, req *http.Request) {
		stripe.Key = os.Getenv("STRIPE_KEY")

		urlParams := req.URL.Query()
		sessIds, ok := urlParams["sessId"]
		if !ok {
			io.WriteString(w, "Invalid Query. Missing \"sessId\"")
			return
		}
		sessionID := sessIds[0]

		s, _ := session.Get(sessionID, nil)

		userName := req.Context().Value(auth.ContextKeyUserName).(string)

		user, _ := storage.Store.User.Get(userName)
		user.IsPro = true
		user.Billing = &storage.UserBillding{
			SubscriptionID: s.Subscription.ID,
		}

		storage.Store.User.Set(userName, user)

		jsonUser, _ := json.Marshal(userName)

		fmt.Printf("%s", jsonUser)

		// io.WriteString(w, "Done!"+s.Subscription.ID)
		http.Redirect(w, req, os.Getenv("CLIENT_URL")+"/pro/success", http.StatusTemporaryRedirect)
	}))

	http.HandleFunc("/payment", auth.WithCredentials(func(w http.ResponseWriter, req *http.Request) {
		origin := req.Header.Get("Origin")
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Content-Type", "application/json")
		stripe.Key = os.Getenv("STRIPE_KEY")
		checkParams := &stripe.CheckoutSessionParams{
			PaymentMethodTypes: stripe.StringSlice([]string{
				"card",
			}),
			// Customer: &customer.ID,
			LineItems: []*stripe.CheckoutSessionLineItemParams{
				{
					Price:    stripe.String("price_1GptsEK2tzGfLmpdXZGBLXjD"),
					Quantity: stripe.Int64(1),
				},
			},
			Mode:       stripe.String("subscription"),
			SuccessURL: stripe.String(os.Getenv("SERVER_URL") + "/payment/success?sessId={CHECKOUT_SESSION_ID}"),
			CancelURL:  stripe.String(os.Getenv("CLIENT_URL") + "/pro"),
		}

		session, err := session.New(checkParams)
		if err != nil {
			w.Write([]byte(err.Error()))
		}

		jsonResp := struct {
			Session string `json:"session"`
		}{
			Session: session.ID,
		}

		res, _ := json.Marshal(jsonResp)

		w.Write(res)
	}))
	return false
}
