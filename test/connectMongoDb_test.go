package test

import (
	"github.com/emtabb/espace"
	"github.com/emtabb/espace/surface"
	"github.com/emtabb/qugo"
	"github.com/emtabb/state"
	"log"
	"testing"
)

func LogState(state state.State) {
	log.Println(state.ToString())
}

func TestLoadSpaceByMongoDb(t *testing.T) {
	var someSpace espace.ESpace = surface.Surface()


	//uri := ""
	//someSpace.MongoSpace(uri, "").LoadSpace("")

	uri := "mongodb://localhost:27017/mongodb"
	someSpace.MongoSpace(uri, "mongodb").LoadSpace("test")


	quantization := qugo.Operator().Init(someSpace.Head())
	quantization.Pipe(LogState)
}