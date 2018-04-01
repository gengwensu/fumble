package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestPOSTInpackets(t *testing.T) { // POST inPackets
	var tests = []struct {
		payload InPacket
		want    int
	}{
		{ //good data
			InPacket{1, 70, -35},
			http.StatusOK,
		},
		{ //bad userId
			InPacket{-1, 70, -35},
			http.StatusBadRequest,
		},
		{ //bad longitude
			InPacket{2, 210, -35},
			http.StatusBadRequest,
		},
		{ //bad longitude
			InPacket{2, -210, -35},
			http.StatusBadRequest,
		},
		{ //bad latitude
			InPacket{3, 70, 110},
			http.StatusBadRequest,
		},
		{ //bad latitude
			InPacket{3, 70, -110},
			http.StatusBadRequest,
		},
	}
	locDb := map[int]Location{}
	handler := AppContext{
		LocDb: locDb,
	}

	server := httptest.NewServer(&handler)
	defer server.Close()

	for _, test := range tests {
		jsonString, err := json.Marshal(test.payload)
		buf := strings.NewReader(string(jsonString))
		url := server.URL + "/fumble/location"
		resp, err := http.Post(url, "application/json", buf)
		if err != nil {
			t.Fatal(err)
		}
		out, err := ioutil.ReadAll(resp.Body)
		if resp.StatusCode == 200 {
			sout := string(out) //convert to string
			swant := "http " + strconv.Itoa(test.want) + "\n"
			if sout != swant {
				t.Errorf("HTTP returns %s, want %s\n", sout, swant)
			}
		} else {
			t.Errorf("HTTPtest server error %d\n", resp.StatusCode)
		}
	}
}

type JsonOut struct {
	key []OutPacket `json:"friends"`
}

func TestGETListFriends(t *testing.T) { // GET - list bundles
	var tests = []struct {
		userId int
		want   OutPacket
	}{
		{
			1,
			OutPacket{1, time.Now(), 2},
		},
		{
			2,
			OutPacket{2, time.Now(), 1},
		},
		{
			3,
			OutPacket{},
		},
	}
	locDb := map[int]Location{
		1: Location{70, -35},
		2: Location{70, -35},
		3: Location{80, -35},
	}
	handler := AppContext{
		LocDb: locDb,
	}
	server := httptest.NewServer(&handler)
	defer server.Close()

	for _, test := range tests {
		baseURL := server.URL + "/fumble/friends/"
		url := baseURL + strconv.Itoa(test.userId)
		resp, err := http.Get(url)
		if err != nil {
			t.Fatal(err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("HTTP returns status code %d; want %d\n", resp.StatusCode, http.StatusOK)
		}

		result, err := ioutil.ReadAll(resp.Body)
		sout := string(result)
		out := OutPacket{}
		if len(sout) > 16 {
			sout = sout[14 : len(sout)-3]
			err = json.Unmarshal([]byte(sout), &out)
			//fmt.Printf("sout %s, out %v\n", sout, out)
			if out.Friend != test.want.Friend ||
				out.User != test.want.User {
				t.Errorf("result %v incorrect, want %v\n", out, test.want)
			}
		} else {
			if sout != "{\n friends: []\n}" {
				t.Errorf("result %s incorrect, want %v\n", sout, "{\n friends: []\n}")
			}
		}
	}
}
