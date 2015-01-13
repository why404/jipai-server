package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"jipai/pili-sdk"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"net/http"
	"net/url"
	"time"
)

type Video struct {
	ID          bson.ObjectId     `bson:"_id" json:"id"`
	StreamId    string            `bson:"stream_id" json:"-"`
	Name        string            `bson:"name" json:"name"`
	Description string            `bson:"description" json:"description"`
	StreamKey   string            `bson:"stream_key" json:"-"`
	PushUrl     string            `bson:"push_url" json:"push_url"`
	LiveUrl     map[string]string `bson:"live_url" json:"live_url"`
	CreatedAt   time.Time         `bson:"created_at" json:"created_at"`

	nonce int
}

func (v *Video) SignPushUrl() string {
	if v.nonce == 0 {
		v.nonce = int(time.Now().Unix())
	} else {
		v.nonce++
	}
	ret := fmt.Sprintf("%s?nonce=%d", v.PushUrl, v.nonce)
	hash := hmac.New(sha1.New, []byte(v.StreamKey))
	hash.Write([]byte(ret))
	token := base64.URLEncoding.EncodeToString(hash.Sum(nil))
	ret = fmt.Sprintf("%s&token=%v", ret, token)
	return ret
}

type Videos struct {
	collection *mgo.Collection
	app        string
	stream     *pili.Streams
	callback   ErrorCallback
}

func NewVideos(mdb *mgo.Database, app, ak, sk string, callback ErrorCallback) (*Videos, error) {
	collection := mdb.C("videos")
	mac := pili.Mac{
		AccessKey: ak,
		SecretKey: sk,
	}
	stream := pili.New(mac)
	return &Videos{
		collection: collection,
		app:        app,
		stream:     stream,
		callback:   callback,
	}, nil
}

func (v *Videos) Create(video Video) (*Video, error) {
	stream, err := v.stream.Create(v.app, "", false)
	if err != nil {
		return nil, Error{http.StatusInternalServerError, err.Error(), 500001}
	}
	video.ID = bson.NewObjectId()
	video.StreamId = stream.ID
	video.StreamKey = stream.StreamKey
	video.PushUrl = stream.PushUrl[0]["RTMP"]
	video.LiveUrl = stream.LiveUrl[pili.StreamDefaultLive]
	video.CreatedAt = time.Now().UTC()
	if err := v.collection.Insert(video); err != nil {
		v.callback.OnError(err)
		return nil, Error{http.StatusInternalServerError, err.Error(), 500001}
	}
	video.PushUrl = video.SignPushUrl()
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
	u, err := url.Parse(video.LiveUrl["RTMP"])
	if err == nil {
		u.Host = "ws1.src.rtmp.pili.qiniu.com"
		video.LiveUrl["RTMP"] = u.String()
	}
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
	ret := struct {
		Url string `json:"push_url"`
	}{}
	ret.Url = video.SignPushUrl()
	return ret, nil
}
