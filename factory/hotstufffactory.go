package factory

import (
	"github.com/wjbbig/go-hotstuff/hotstuff"
	"github.com/wjbbig/go-hotstuff/hotstuff/basic"
)

func HotStuffFactory(networkType string, id int) hotstuff.HotStuff {
	switch networkType {
	case "basic":
		return basic.NewBasicHotStuff(id)
	case "chained":
		break
	case "event-driven":
		break
	}
	return nil
}

