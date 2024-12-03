package types

import (
	"time"
)

type Request struct {
	URL         string        `json:"url"`
	CustomShort string        `json:"customShort"`
	ExpiryDate  time.Duration `json:"expiryDate"`
}
type Response struct {
	URL            string        `json:"url"`
	CustomShort    string        `json:"customShort"`
	Expiry         time.Duration `json:"expiry"`
	XRateRemaining int           `json:"xRateRemaining"`
	XRateLimitRest time.Duration `json:"xRateLimitRest"`
}

type QrResponse struct {
	URL            string        `json:"url"`
	QrCode         []byte        `json:"qrCode"`  // Change from string to []byte for raw binary data
	Expiry         time.Duration `json:"expiry"`
	XRateRemaining int           `json:"xRateRemaining"`
	XRateLimitRest time.Duration `json:"xRateLimitRest"`
}
