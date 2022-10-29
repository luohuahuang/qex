package config

const (
	GitMasterZooKeeper  = "example.com:2181"
	GitMasterBootstramp = "example.com:9092"
	GitMasterReport     = "GIT_MASTER_LOG"
	GitMasterTopic      = "qex-git-topic"

	JenkinsZooKeeper          = "example.com:2181"
	JenkinsBootstramp         = "example.com:9092"
	JenkinsBuildTopic         = "qex-jenkins-build-topic"
	JenkinsBuildExporterTopic = "qex-jenkins-build-exporter-topic"

	JiraServer = "https://jira.example.com/"

	InfluxBucket = "example.com_test_report_db"
	InfluxOrg    = "example.com-qa"
	InfluxToken  = "token-not-required"
	InfluxUrl    = "http://example.com:8086"

	CacheServer = "example.com:6379/0"

	JenkinsBuildInfoURL = "http://platform.qa.example.com/automation/job/%s/%s/api/json"

	GitReadOnlyToken = "provide-your-token-here"
	GitV4API         = "https://example.com/api/v4"
)
