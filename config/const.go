package config

const (
	BaseDir = "/tmp/"

	TestExecutionMeasurement = "all_records_v2"

	GitMasterLogFormat         = "git-master-log-%s.txt"
	GitMasterLogFullPathFormat = BaseDir + GitMasterLogFormat
	GitMeasurement             = "all_git_records_v2"

	JiraOKRMeasurement = "all_jira_okr_records_v2"

	GitMRMeasurement   = "all_git_mr_records_v2"
	GitUserCacheFormat = "qex-user-%d"

	JenkinsBuildMeasurement = "all_jenkins_build_records_v2"
)

var (
	MapATSignOff = map[string]string{
		"product-1": "labels in (xxx)",
		"product-2": "labels in (xxx)",
	}

	MapATFoundBug = map[string]string{
		"product-1": "labels in (xx)",
		"product-2": "labels in (xxx)",
	}

	MapGitTestRepo = map[string]int{
		"product-1": 38753,
		"product-2": 36142,
	}

	MapGitProductRepo = map[string]int{
		"huanglh/my-abc-service": 57841,
	}
)
