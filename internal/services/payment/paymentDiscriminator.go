package payment

import (
	"errors"
	"icecreamshop/internal/messageErrors"
)

// ProcessPayment discriminates payment data and process the payment method chosen
func ProcessPayment(paymentData PaymentRequest, totalCost uint) (any, error) {
	switch paymentData.PaymentType {
	case CreditCardType:
		if paymentData.CreditCard == nil {
			return "", errors.New(messageErrors.InvalidPaymentData)
		}
		paymentData.CreditCard.Amount = totalCost
		return paymentData.CreditCard.Process()
	case DigitalWalletType:
		if paymentData.DigitalWallet == nil {
			return "", errors.New(messageErrors.InvalidPaymentData)
		}
		paymentData.DigitalWallet.Amount = totalCost
		return paymentData.DigitalWallet.Process()
	case PreferenceMPType:
		if paymentData.PreferenceMP == nil {
			return "", errors.New(messageErrors.InvalidPaymentData)
		}
		paymentData.PreferenceMP.Amount = totalCost
		return paymentData.PreferenceMP.Process()
	default:
		return "", errors.New(messageErrors.UnsupportedPaymentType)
	}
}
