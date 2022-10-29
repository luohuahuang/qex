package config

const (
	BaseDir = "/tmp/"

	XmlFileDir = "/home/user/.jenkins/workspace/"

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
	}

	MapATFoundBug = map[string]string{
	}

	MapGitTestRepo = map[string]int{
		"example":       38753,
	}

	MapGitProductRepo = map[string]int{
		"huanglh/my-abc-service/": 57841,
	}

	MatterMostPublic  = "https://mattermost.example.com/hooks/ptxk18boipn4ibibxw1758gueh" // qex-msg-bot
	MatterMostMonitor = "https://mattermost.example.com/hooks/8mmogp5jo7rrfe6kprgwf45mdo" // qex-service-alive-bot

	MatterMostEndpoints = map[string]string{
		"^qex-": "https://mattermost.example.com/hooks/8mmogp5jo7rrfe6kprgwf45mdo", // qex-service-alive-bot
	}
)
