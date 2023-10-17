package test

import (
	"log"
	"testing"

	"github.com/anacrolix/torrent"
)

func TestExample(t *testing.T) {
	c, err := torrent.NewClient(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	tor, err := c.AddMagnet("magnet:?xt=urn:btih:d0d2d1482fd98ab4650879fdca625da57fbeedb0&tr=http://open.acgtracker.com:1096/announce")
	if err != nil {
		t.Fatal(err)
	}

	<-tor.GotInfo()
	tor.DownloadAll()
	c.WaitAll()
	log.Print("ermahgerd, torrent downloaded")
}
