package config

const (
	BaseDir = "/tmp/"

	TestExecutionMeasurement = "all_records_v2"

	GitMasterLogFormat         = "git-master-log-%s.txt"
	GitMasterLogFullPathFormat = BaseDir + GitMasterLogFormat
	GitMeasurement             = "all_git_records_v2"

	JiraOKRMeasurement = "all_jira_okr_records_v2"

	GitMasterReport = "GIT_MASTER_LOG"
	GitMasterTopic  = "qex-git-topic"
)
