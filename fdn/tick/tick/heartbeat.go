package tick

import (
	"github.com/hootuu/gelato/errors"
)

func Heartbeat(id ID) *errors.Error {
	//todo
	//err := mgo.TxSyncModify(context.Background(), tickNativeColl(), bson.M{
	//	"id": id,
	//}, bson.M{
	//	"$set": bson.M{
	//		"lst_heartbeat_time": time.Now(),
	//		"updated_at":        time.Now(),
	//	},
	//	"$inc": bson.M{
	//		"version_id": 1,
	//	},
	//})
	//if err != nil {
	//	return err
	//}

	return nil
}
