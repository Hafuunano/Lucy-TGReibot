package blocker

import (
	ctrl "github.com/FloatTech/zbpctrl"
	rei "github.com/fumiama/ReiBot"
)

var (
	engine = rei.Register("blocker", &ctrl.Options[*rei.Ctx]{
		DisableOnDefault: true,
		Help:             "block",
	})
)

func init() {
}
