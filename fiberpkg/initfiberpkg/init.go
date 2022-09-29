package initfiberpkg

import "github.com/gofiber/fiber/v2"

type InitialAppStruct struct {
	AllowOrigin      []string
	AllowMethod      string
	AllowHeader      string
	AllowCredentials bool
	MaxAge           int
	Routes           []Route
	PrefixPath       string
	State            string
	Path             string
}

type Route struct {
	Name       string
	Method     string
	Pattern    string
	Endpoint   fiber.Handler
	Middleware []fiber.Handler
	Test       bool
}
