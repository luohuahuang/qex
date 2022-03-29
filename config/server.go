package config

const (
	GitMasterZooKeeper  = "example.com:2181" //
	GitMasterBootstramp = "example.com:9092"
	GitMasterReport     = "GIT_MASTER_LOG"
	GitMasterTopic      = "qex-git-topic"

	JenkinsZooKeeper          = "example.com:2181"
	JenkinsBootstramp         = "example.com:9092"
	JenkinsBuildTopic         = "qex-jenkins-build-topic"
	JenkinsBuildExporterTopic = "qex-jenkins-build-exporter-topic"

	JiraServer = "https://jira.example.com/"

	InfluxBucket = "example_test_report_db"
	InfluxOrg    = "example-qa"
	InfluxToken  = "token-not-required"
	InfluxUrl    = "http://ywxsji.vm.cloud.example.com:8086"

	MatterMost = "https://mattermost.example.com/hooks/5t87430jgoerjgoerjt"

	CacheServer = "1agmh5.vm.cloud.example.com:6379/8"

	JenkinsBuildInfoURL = "http://1xh5ym.vm.cloud.example.com:8080/job/%s/%s/api/json"

	GitReadOnlyToken = "g2_yocXCFhAR7cxd2gn2"
	GitV4API         = "https://sec3.git.garena.com/api/v4"
)
