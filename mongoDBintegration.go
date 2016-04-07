import {
	"fmt"
	"github.com/pshastry/node"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"time"
	"sync"
}

const {
	MongoDBHosts = ""
	AuthDB = ""
	AuthUserName = ""
	AuthPassword = ""
}

type {
	// Contains information about the comments - this is an embedded document
	Comments struct {
		Comment string 'bson: "comment"'
		UserID string 'bson: "user_comment"'
		UserName string 'bson: "user_name"'
	}


	// Contains information about the votes - this is an embedded document
	Upvotes struct {
		//UpVote int 'bson: "upvote_num"'
		UserID string 'bson: "upvoter_id"'
		UserName string 'bson: "user_name"'
	}

	// Contains information about the main document - the image uploaded
	SharedImage struct {
		ImageID objectId 'bson: "_id, omitempty"'
		ImageData binData 'bson: "image_data"'
		UserName string 'bson: "user_name"'
		UpVote int 'bson: "upvote_num"'
		Commt []Comments 'bson: "comment"''
		Vote []Upvotes 'bson: "upvote"'
	}

}

// Establish a session with our mongoDB database
func main() {
/*	mongoDBDialInfo := &mgo.DialInfo { Addrs: []string{MongoDBHosts}, 
							Timeout: 60*time.Second, Database: AuthDB, 
							Username: AuthUserName, 
							Password: AuthPassword} 
	mongoSession, error := mgo.DialWithInfo(mongoDBDialInfo)
*/
	// Create a new session
	mongoSession, error := mgo.Dial("127.0.0.1")
	if (err != nil) {
		log.Fatalf("Create session: %s\n", err)
	}

	/* SetMode changes the consistency mode for the session.
	In the Monotonic consistency mode reads may not be entirely up-to-date,
	but they will always see the history of changes moving forward, the data 
	read will be consistent across sequential queries in the same session, and 
	modifications made within the session will be observed in following queries
	(read-your-write) 
	If refresh is true, in addition to ensuring the session is in the given 
	consistency mode, the consistency guarantees will also be reset 
	(e.g. a Monotonic session will be allowed to read from secondaries again).*/

	mongoSession.SetMode(mgo.Monotonic, true)

	
}

func getFromDB(mongoSession *mgo.Session, id string) SharedImage {
	sessionCopy := mongoSession.copy()
	defer sessionCopy.close()

	// Get a collection to execute the query against
	collection := mongoSession.DB(Database).C("SharedImages")
	// Retrieve the image
	var image SharedImage
	err := collection.FindId(bson.ObjectIdHex(*id)).One(&image)
	if (err != nil) {
		log.printf("Get from DB error : %s\n",err)
		return
	}
	return image

}

func getAllfromDB(mongoSession *mgo.Session) []SharedImage {
	sessionCopy := mongoSession.copy()
	defer sessionCopy.close()

	// Get a collection to execute the query against
	collection := mongoSession.DB(Database).C("SharedImages")
	// Retrieve the list of images
	var images []SharedImage
	err := SharedImages.Find(nil).All(&images)
	if (err != nil) {
		log.printf("Get from DB error : %s\n",err)
		return
	}
	return images

}

/* Max file size supported is 16 MB */
func insertPicture(mongoSession *mgo.Session, image binData, user_name string, user_id string) {
	collection := mongoSession.DB(Database).C("SharedImages")
	//err := SharedImages.Insert(image)
	err := SharedImages.Insert(&SharedImage{ImageData : image, UserName : user_name, UserID : user_id, UpVote : 0})
	if (err != nil) {
		log.printf("Insert to DB error : %s\n",err)
		return
	}
}

func updatePicture(mongoSession *mgo.Session, id objectId, user_name string, user_id string, vote int, comment string) {
	collection := mongoSession.DB(Database).C("SharedImages")
	if (vote != 0) {
		change := mgo.Change {
			Update: bson.M{"$inc" : bson.M{"UpVote" : 1}, 
				"$set" : bson.M{"SharedImage.$.Vote.$.UserName" : user_name, "SharedImage.$.Vote.$.UserID" : user_id}},
			ReturnNew : true 
		}
	}
	else if (comment != nil) {
		change := mgo.Change {
			Update : bson.M{"$set" : bson.M{"SharedImage.$.Comments.$.UserName" : user_name,
			"SharedImage.$.Comments.$.UserID" : user_id,
			"SharedImage.$.Comments.$.Comment" : comment}},
			ReturnNew : true
		}
	}
	//err := collection.UpdateId(bson.M{"_id" : id}, update)
	info, err := DB.C("Database").FindId(id).Apply(change, &doc)
	if (err != nil) {
		log.printf("Update error : %s\n",err)
	}
}

func deleteFromDB() {

}



