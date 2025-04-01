package handlers

type stub struct {
	Title, CssFile string
	LoggedIn       bool
}

var stubs = map[string]stub{
	"about":           {Title: "About"},
	"contact":         {Title: "Contact"},
	"faq":             {Title: "FAQ"},
	"features":        {Title: "Features"},
	"forgot-password": {Title: "Forgot Password"},
	"index":           {Title: "Home", CssFile: "index.css"},
	"login":           {Title: "Login"},
	"profile":         {Title: "My Profile"},
	"signup":          {Title: "Sign Up"},
	"pricing":         {Title: "Pricing"},
}
