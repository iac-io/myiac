package elasticsearch

func CreateSnapshotRepository() {
	//curl -i -H "Content-Type: application/json" -XPUT http://localhost:9200/_snapshot/moneycol-banknotes-backup -d '{ "type": "fs", "settings": {"location": "/tmp/backups"}}'
}

func executeCurlCmd() {

}