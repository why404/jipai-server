package main

import (
	"jipai/pili-sdk"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"net/http"
	"time"
)

type Video struct {
	Id          bson.ObjectId `bson:"_id" json:"id"`
	StreamId    string        `bson:"stream_id" json:"-"`
	Name        string        `bson:"name" json:"name"`
	Description string        `bson:"description" json:"description"`
	StreamKey   string        `bson:"stream_key" json:"-"`
	PushUrl     string        `bson:"push_url" json:"push_url,omitempty"`
	LiveUrl     struct {
		HLS  string `bson:"hls" json:"HLS"`
		RTMP string `bson:"rtmp" json:"RTMP"`
	} `bson:"live_url" json:"live_url"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}

type Videos struct {
	collection *mgo.Collection
	stream     *pili.Streams
	callback   ErrorCallback
}

func NewVideos(mdb *mgo.Database, ak, sk string, callback ErrorCallback) (*Videos, error) {
	collection := mdb.C("videos")
	mac := pili.Mac{
		AccessKey: ak,
		SecretKey: sk,
	}
	stream := pili.New(mac)
	return &Videos{
		collection: collection,
		stream:     stream,
		callback:   callback,
	}, nil
}

func (v *Videos) Create(video Video) (*Video, error) {
	stream, err := v.stream.Create(video.Name, "RTMP", "", false, 0)
	if err != nil {
		return nil, Error{http.StatusInternalServerError, err.Error(), 500001}
	}
	video.Id = bson.NewObjectId()
	video.StreamKey = stream.StreamKey
	video.PushUrl = stream.PushUrl
	video.LiveUrl.HLS = stream.LiveUrl.HLS
	video.LiveUrl.RTMP = stream.LiveUrl.RTMP
	video.CreatedAt = time.Now().UTC()
	if err := v.collection.Insert(video); err != nil {
		v.callback.OnError(err)
		return nil, Error{http.StatusInternalServerError, err.Error(), 500001}
	}
	video.PushUrl = stream.SignPushUrl()
	return &video, nil
}

func (v *Videos) List() ([]Video, error) {
	var ret []Video
	if err := v.collection.Find(bson.M{}).All(&ret); err != nil {
		v.callback.OnError(err)
		return nil, Error{http.StatusInternalServerError, err.Error(), 500001}
	}
	for i := range ret {
		ret[i].PushUrl = ""
		ret[i].CreatedAt = ret[i].CreatedAt.UTC()
	}
	return ret, nil
}

func (v *Videos) Get(id string) (*Video, error) {
	var video Video
	if err := v.collection.FindId(bson.ObjectIdHex(id)).One(&video); err != nil {
		if err == mgo.ErrNotFound {
			return nil, Error{http.StatusNotFound, "not found", 401001}
		}
		v.callback.OnError(err)
		return nil, Error{http.StatusInternalServerError, err.Error(), 500001}
	}
	video.PushUrl = ""
	video.CreatedAt = video.CreatedAt.UTC()
	return &video, nil
}

func (v *Videos) Delete(id string) error {
	video, err := v.Get(id)
	if err != nil {
		return err
	}
	if _, err := v.stream.Delete(video.StreamId); err != nil {
		return Error{http.StatusInternalServerError, err.Error(), 500001}
	}
	if err := v.collection.RemoveId(bson.ObjectIdHex(id)); err != nil {
		if err == mgo.ErrNotFound {
			return Error{http.StatusNotFound, "not found", 401001}
		}
		v.callback.OnError(err)
		return Error{http.StatusInternalServerError, err.Error(), 500001}
	}
	return nil
}

func (v *Videos) GetPushUrl(id string) (interface{}, error) {
	video, err := v.Get(id)
	if err != nil {
		return nil, err
	}
	stream := pili.Stream{
		PushUrl:   video.PushUrl,
		StreamKey: video.StreamKey,
	}
	ret := struct {
		Url string `json:"push_url"`
	}{}
	ret.Url = stream.SignPushUrl()
	return ret, nil
}
