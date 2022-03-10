package config

const (
	// UPDATE THE CONFIG WITH YOUR SETTING - START HERE
	GitMasterZooKeeper  = "kafka.example.com:2181"
	GitMasterBootstramp = "kafka.example.com:9092"
	JiraServer          = "https://jira.example.com/"

	InfluxUrl    = "http://influx.example.com:8086"
	InfluxBucket = "qex_db"
	InfluxOrg    = "qa-organization"
	InfluxToken  = "fill-in-token-if-needed"

	MsgServer = "https://mattermost.example.com/hooks/ogrnuzhsej8sifrgkgb1scd96y"
	// UPDATE THE CONFIG WITH YOUR SETTING - END HERE
)

var (
	MapATSignOff = map[string]string{
		"product-line-1": "labels in (product-1-AT-sign-off, product-1-AT-sign-off-fully, product-1-AT-sign-off-partial)",
		"product-line-2": "labels in (product-2-AT-sign-off, product-2-AT-sign-off-fully, product-2-AT-sign-off-partial)",
		// more product liens here, ask your product lines for the JQL they define for the AT sign off
	}

	MapATFoundBug = map[string]string{
		"product-line-1": "labels in (product-1-AT-found, product-1-AT-found-nightly, product-1-AT-found-sign-off)",
		"product-line-2": "labels in (product-2-AT-Found, product-2-AT-Found-Nightly, product-2-AT-Found-Regression, product-2-AT-found-nightly, product-2-AT-found-bug-nightly)",
		// more product liens here, ask your product lines for the JQL they define for the AT found bug
	}
)
