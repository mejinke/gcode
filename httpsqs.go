package gcode

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type Httpsqs struct {
	host, port, auth, charset string
}

//Init a new httpsqs
func NewHttpSQS(options ...string) *Httpsqs {
	mq := &Httpsqs{"localhost", "1218", "", "utf-8"}
	for i := 0; i < len(options); i++ {
		switch i {
		case 0:
			mq.host = options[0]
		case 1:
			mq.port = options[1]
		case 2:
			mq.auth = options[2]
		case 3:
			mq.charset = options[3]
		}
	}
	return mq
}

//Build query string for httpsqs
func (mq *Httpsqs) makeQuery(query string) (res string) {
	query = "http://" + mq.host + ":" + mq.port + "/?" +
		"auth=" + mq.auth + "&charset=" + mq.charset + "&" + query
	return query
}

//Do Put request from httpsqs
//func (mq *Httpsqs) put(query string, value string) (res string, err error) {
//}

//put data to queue
func (mq *Httpsqs) Put(queue string, value string) (rs bool, err error) {
	query := mq.makeQuery("name=" + queue + "&opt=put")
	r := new(http.Response)
	r, err = http.Post(query, "", strings.NewReader(value))
	if err != nil {
		return
	}
	defer r.Body.Close()
	rBytes := []byte{}
	rBytes, err = ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	res := string(rBytes)
	if res == "HTTPSQS_PUT_OK" {
		return true, nil
	}
	return false, errors.New(res)
}

//Do Get request from httpsqs
func (mq *Httpsqs) get(query string) (res string, err error) {
	r := new(http.Response)
	r, err = http.Get(mq.makeQuery(query))
	if err != nil {
		return
	}
	defer r.Body.Close()
	rBytes := []byte{}
	rBytes, err = ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	res = string(rBytes)
	return
}

//Do Get request from httpsqs
//And gets normal string data
func (mq *Httpsqs) getString(query string) (res string, err error) {
	res, err = mq.get(query)
	if err != nil {
		return "", err
	}
	if res == "HTTPSQS_ERROR" {
		return "", errors.New(res)
	}
	return
}

//Do Get request from httpsqs
//And gets normal bool data
func (mq *Httpsqs) getBool(query string, expected string) (rs bool, err error) {
	var res string
	res, err = mq.get(query)
	if err != nil {
		return false, err
	}
	return res == expected, nil
}

//Get data from queue with position
func (mq *Httpsqs) PGet(queue string) (res string, pos int, err error) {
	query := mq.makeQuery("name=" + queue + "&opt=get")
	r := new(http.Response)
	r, err = http.Get(query)
	if err != nil {
		return
	}
	defer r.Body.Close()
	rBytes := []byte{}
	rBytes, err = ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	res = string(rBytes)
	if res == "HTTPSQS_ERROR" {
		return "", 0, errors.New(res)
	}
	if res != "HTTPSQS_GET_END" {
		var posTmp string
		posTmp = r.Header.Get("pos")
		if posTmp != "" {
			pos, err = strconv.Atoi(posTmp)
		}
	} else {
		res = ""
	}
	return res, pos, nil
}

//Get data from queue
func (mq *Httpsqs) Get(queue string) (res string, err error) {
	query := "name=" + queue + "&opt=get"
	res, err = mq.getString(query)
	if err != nil {
		return "", err
	}
	if res == "HTTPSQS_GET_END" {
		res = ""
	}
	return
}

//Get status from queue
func (mq *Httpsqs) Status(queue string) (res string, err error) {
	query := "name=" + queue + "&opt=status"
	return mq.getString(query)
}

//Get status from queue in json format
func (mq *Httpsqs) StatusJson(queue string) (res string, err error) {
	query := "name=" + queue + "&opt=status_json"
	return mq.getString(query)
}

//View data from queue
func (mq *Httpsqs) View(queue string, pos int) (res string, err error) {
	query := "name=" + queue + "opt=view&pos=" + string(pos)
	return mq.getString(query)
}

//Clear queue
func (mq *Httpsqs) Reset(queue string) (rs bool, err error) {
	query := "name=" + queue + "&opt=reset"
	return mq.getBool(query, "HTTPSQS_RESET_OK")
}

//Modify the maximum of queue
func (mq *Httpsqs) MaxQueue(queue string, num int) (rs bool, err error) {
	query := "name=" + queue + "&opt=maxqueue&num=" + string(num)
	return mq.getBool(query, "HTTPSQS_MAXQUEUE_OK")
}

//Modify the frequecy for httpsqs to save data to disk
func (mq *Httpsqs) SyncTime(num int) (rs bool, err error) {
	query := "name=httpsqs_synctime&opt=synctime&num=" + string(num)
	return mq.getBool(query, "HTTPSQS_SYNCTIME_OK")
}
