package payment

import (
	"context"
	"errors"
	"fmt"
	"icecreamshop/internal/messageErrors"
	"os"

	"github.com/mercadopago/sdk-go/pkg/config"
	"github.com/mercadopago/sdk-go/pkg/preference"
)

const CreditCardType = "creditCard"
const DigitalWalletType = "digitalWallet"
const PreferenceMPType = "preferenceMP"

type Payment struct {
	Amount      uint   `json:"amount"`
	PaymentType string `json:"payment_type"`
}

type CreditCard struct {
	Payment
	CardNumber      string `json:"card_number"`
	ExpirationMonth string `json:"expiration_month"`
	ExpirationYear  string `json:"expiration_year"`
	CVV             string `json:"cvv"`
	CardHolderName  string `json:"card_holder_name"`
}

type DigitalWallet struct {
	Payment
	WalletID string `json:"wallet_id"`
}

type PreferenceMP struct {
	Payment
}

type PaymentRequest struct {
	PaymentType   string         `json:"payment_type"`
	DigitalWallet *DigitalWallet `json:"digital_wallet,omitempty"`
	CreditCard    *CreditCard    `json:"credit_card,omitempty"`
	PreferenceMP  *PreferenceMP  `json:"preference_mp,omitempty"`
}

func (p *CreditCard) Process() (string, error) {
	err := p.Validate()
	if err != nil {
		return "", err
	}

	//External integration service is needed
	//For testing, we assume it paid.

	return "ABCDE123", nil
}

func (p *PreferenceMP) Process() (*preference.Response, error) {
	if p.Amount <= 0 {
		return &preference.Response{}, errors.New(messageErrors.MustBeAnInteger)
	}

	accessToken := os.Getenv("MP_ACCESS_TOKEN")
	cfg, err := config.New(accessToken)
	if err != nil {
		fmt.Println(err)
	}

	client := preference.NewClient(cfg)

	request := preference.Request{
		Items: []preference.ItemRequest{
			{
				Title:     "icrecreamshop-backend Payment",
				Quantity:  1,
				UnitPrice: float64(p.Amount),
			},
		},
	}

	resource, err := client.Create(context.Background(), request)
	if err != nil {
		return &preference.Response{}, err
	}

	return resource, nil
}

func (p *DigitalWallet) Process() (string, error) {
	err := p.Validate()
	if err != nil {
		return "", err
	}

	//External integration service is needed
	//For testing, we assume it paid.

	return "ABCDE123", nil
}

func (p *CreditCard) Validate() error {
	if p.Amount <= 0 {
		return errors.New(messageErrors.MustBeAnInteger)
	}
	if len(p.CardNumber) != 16 {
		return errors.New(messageErrors.InvalidCreditCardLength)
	}
	if len(p.CardHolderName) == 0 {
		return errors.New(messageErrors.CreditCardHolderNameIsRequired)
	}
	if len(p.ExpirationMonth) != 2 {
		return errors.New(messageErrors.InvalidExpirationMonth)
	}
	if len(p.ExpirationYear) != 4 {
		return errors.New(messageErrors.InvalidExpirationYear)
	}
	if len(p.CVV) != 3 {
		return errors.New(messageErrors.InvalidCVV)
	}
	return nil
}

func (p *DigitalWallet) Validate() error {
	if p.Amount <= 0 {
		return errors.New(messageErrors.MustBeAnInteger)
	}
	if len(p.WalletID) == 0 {
		return errors.New(messageErrors.WalletIDIsRequired)
	}
	return nil
}
