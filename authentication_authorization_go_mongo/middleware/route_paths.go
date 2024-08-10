package middleware

var PathMap = map[string]map[string]bool{
	"": map[string]bool{
		"user":  true,
		"admin": true,
	},
	"GET": map[string]bool{
		"user":  true,
		"admin": true,
	},
	"POST": map[string]bool{
		"user":  false,
		"admin": true,
	},
	"PUT": map[string]bool{
		"user":  false,
		"admin": true,
	},
	"DELETE": map[string]bool{
		"user":  false,
		"admin": true,
	},
	"PATCH": map[string]bool{
		"user":  false,
		"admin": true,
	},
}
