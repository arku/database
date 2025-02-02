package models

import (
	"testing"
)

func TestRecord(t *testing.T) {
	var rm = NewRecordManager(newTestDB(t, &Record{}))
	type args struct {
		username      string
		recordName    string
		recordKeyName string
		zoneName      string
		ipfsHash      string
		metadata      map[string]interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{"NoMetaData", args{"testuser1", "testrecord1", "testkey1", "testzone1", "testhash1", nil}},
		{"YesMetaData", args{"testuser2", "testrecord2", "testkey2", "testzone2", "testhash2", map[string]interface{}{
			"food": "pizza",
			"pet":  "dog",
		}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			record1, err := rm.AddRecord(
				tt.args.username,
				tt.args.recordName,
				tt.args.recordKeyName,
				tt.args.zoneName,
				tt.args.metadata,
			)
			if err != nil {
				t.Fatal(err)
			}
			defer rm.DB.Delete(record1)
			if record1.LatestIPFSHash != "" {
				t.Fatal("latest ipfs hash should be empty")
			}
			record2, err := rm.UpdateLatestIPFSHash(
				tt.args.username,
				tt.args.recordName,
				tt.args.ipfsHash,
			)
			if err != nil {
				t.Fatal(err)
			}
			if record2.LatestIPFSHash != tt.args.ipfsHash {
				t.Fatal("bad ipfs hash set")
			}
			record3, err := rm.FindRecordByNameAndUser(tt.args.username, tt.args.recordName)
			if err != nil {
				t.Fatal(err)
			}
			if record3.LatestIPFSHash != tt.args.ipfsHash {
				t.Fatal("bad record recovered")
			}
			records, err := rm.FindRecordsByZone(tt.args.username, tt.args.zoneName)
			if err != nil {
				t.Fatal(err)
			}
			if len(*records) != 1 {
				t.Fatal("bad amount of records recovered")
			}
		})
	}
}
