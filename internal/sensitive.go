package internal

var sensitives = map[string]struct{}{
	"password":      {},
	"passwd":        {},
	"pass":          {},
	"credential":    {},
	"secret":        {},
	"token":         {},
	"apikey":        {},
	"signature":     {},
	"authorization": {},
	"email":         {},
	"phone":         {},
}
