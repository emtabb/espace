package test

import (
	"github.com/emtabb/espace"
	. "github.com/emtabb/espace/api/element"
	"github.com/emtabb/espace/surface"
	"github.com/emtabb/qugo"
	"github.com/emtabb/state"
	"log"
	"testing"
)

func LogState(state state.State) {
	log.Println(state.(*Element).ToString())
}

func TestLoadSpaceByMongoDb(t *testing.T) {
	var someSpace espace.ESpace = surface.Surface()
	//uri := ""
	//someSpace.MongoSpace(uri, "").LoadSpace("")

	uri := "mongodb://localhost:27017/izanami"
	someSpace.MongoSpace(uri, "izanami").LoadSpace("page")
	quantization := qugo.Operator().InitStates(new(state.List).ByStates(someSpace.Head()))
	quantization.Pipe(LogState)
}