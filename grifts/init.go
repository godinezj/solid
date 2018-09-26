package grifts

import (
	"bitbucket.org/godinezj/solid/actions"
	"github.com/gobuffalo/buffalo"
)

func init() {
	buffalo.Grifts(actions.App())
}
