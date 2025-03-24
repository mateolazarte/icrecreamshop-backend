package messageErrors

const (
	//User messageErrors
	UserIDNotFound         = "No user found with this ID."
	UserEmailNotFound      = "No user found with this email."
	EmailAlreadyExists     = "Email already exists."
	UserIsAlreadyAnAdmin   = "User is already an admin."
	UserIsAlreadyADriver   = "User is already a delivery driver."
	InvalidEmailOrPassword = "Invalid email or password."
	EmailIsRequired        = "Email is required."
	FirstNameIsRequired    = "First name is required."
	LastNameIsRequired     = "Last name is required."
	PasswordIsTooShort     = "Password must be at least 8 characters long."
	AlreadyLoggedIn        = "Already logged in."

	//General messageErrors
	InvalidJsonFormat           = "Invalid json format."
	ErrorWhileProcessingRequest = "An error occurred while processing your request. Please try again later."
	MustBeAnInteger             = "The value must be a positive number"

	//Order messageErrors
	OrderNotFound          = "No order found with this ID."
	FlavorNotFound         = "No flavor found with this ID."
	AlreadyExistingFlavor  = "This flavor ID already exists."
	NonExistingFlavors     = "One or more flavor do not exist."
	WeightNotAvailable     = "Weight not available."
	WeightCannotBeZero     = "Weight must be a positive number."
	IceCreamTubNotFound    = "No ice cream tub found with this ID."
	FlavorIdIsRequired     = "Flavor id is required."
	FlavorNameIsRequired   = "Flavor name is required."
	FlavorTypeIsRequired   = "Flavor type is required."
	AddressIsRequired      = "Address is required."
	InvalidAmountOfFlavors = "Flavors cannot be 0 or greater than 4."

	//Delivery drivers messageErrors
	DeliveryDriverNotFound      = "No delivery driver found with this ID."
	InvalidCuilFormat           = "Cuil must be 10 or 11 digits long"
	AgeMustBeGreaterThan18      = "Age must be equal or greater than 18."
	AtLeastOneVehicleIsRequired = "You must have at least one vehicle."
	InvalidVehicleIDFormat      = "Vehicle id must be between 6 and 7 digits"

	//Payment messageErrors
	UnsupportedPaymentType         = "Unsupported payment type."
	InvalidPaymentData             = "Invalid payment data."
	WalletIDIsRequired             = "Wallet id is required."
	InvalidCreditCardLength        = "Credit card must have 16 digits."
	CreditCardHolderNameIsRequired = "Credit card holder name is required."
	InvalidExpirationMonth         = "Invalid expiration month."
	InvalidExpirationYear          = "Invalid expiration year."
	InvalidCVV                     = "Invalid CVV."
)
