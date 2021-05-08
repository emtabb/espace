package test

import (
	"github.com/emtabb/espace"
	"github.com/emtabb/espace/surface"
	"github.com/emtabb/qugo"
	"github.com/emtabb/state"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"testing"
)

func LogState(state state.State) {
	log.Println(state)
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

type PageDomainForTest struct {
	state.State
	ID primitive.ObjectID `bson:"_id" json:"_id"`
	Author string `bson:"author" json:"author"`
	Description string `bson:"description" json:"description"`
	SearchTitle string `bson:"search_title" json:"search_title"`
}

func TestLoadModelByMongoDb(t *testing.T) {
	var someSpace espace.ESpace = surface.Surface()
	//uri := ""
	//someSpace.MongoSpace(uri, "").LoadSpace("")

	uri := "mongodb://localhost:27017/izanami"
	var model []PageDomainForTest
	someSpace.MongoSpace(uri, "izanami").LoadSpace("page")
	listSpace := someSpace.LoadModel(&model)

	qugo.Operator().InitStates(listSpace).ForEach(func(s state.State) {
		log.Println(s)
	})
}