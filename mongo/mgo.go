package mongo

import (
	"os"

	mgo "gopkg.in/mgo.v2"
)

var (
	//DB define DB to use
	DB *mgo.Database
)

func init() {
	session, err := mgo.Dial(os.Getenv("MONGO_ADDR"))
	if err != nil {
		panic(err)
	}

	session.SetMode(mgo.Monotonic, true)

	DB = session.DB(os.Getenv("MONGO_DB"))
}

//Close the session
func Close() {
	DB.Session.Close()
}
