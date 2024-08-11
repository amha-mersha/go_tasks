package middleware

var PathMap = map[string]map[string]bool{
	"": {
		"user":  true,
		"admin": true,
	},
	"GET": {
		"user":  true,
		"admin": true,
	},
	"POST": {
		"user":  false,
		"admin": true,
	},
	"PUT": {
		"user":  false,
		"admin": true,
	},
	"DELETE": {
		"user":  false,
		"admin": true,
	},
	"PATCH": {
		"user":  false,
		"admin": true,
	},
}
