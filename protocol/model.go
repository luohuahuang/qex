package protocol

type QEXTestCase struct {
	Product    string  `json:"product"`
	SubProduct string  `json:"sub_product"`
	Service    string  `json:"service"`
	API        string  `json:"api"`
	RunId      string  `json:"run_id"`
	Case       string  `json:"case"`
	Branch     string  `json:"branch"`
	Maintainer string  `json:"maintainer"`
	Timestamp  int64   `json:"timestamp"`
	Duration   float64 `json:"duration"`
	Status     int     `json:"status"`
}

type Cases struct {
	Data []string `json:"data"`
}

type Maintainers struct {
	Data map[string]string `json:"data"`
}

type JenkinsBuild struct {
	JobName    string `json:"job_name"`
	BuildNum   string `json:"build_num"`
	NodeName   string `json:"node_name"`
	BuildCause string `json:"build_cause"`
	GitBranch  string `json:"git_branch"`
	BuildUrl   string `json:"build_url"`
	JiraNum    string `json:"Jira_num"`
	QAPNum     string `json:"qap_num"`
}

type Git struct {
	RunId      string `json:"run_id"`
	CommitId   string `json:"commit_id"`
	Maintainer string `json:"maintainer"`
	Product    string `json:"product"`
	Case       string `json:"case"`
}

const (
	OKRTypeATSignOff  = 0
	OKRTypeATFoundBug = 1
)

type Jira struct {
	RunId   string `json:"run_id"`
	Product string `json:"product"`
	JiraId  string `json:"jira_id"`
	OKRType int    `json:"okr_type"`
}

type GitMR struct {
	RunId   string `json:"run_id"`
	Product string `json:"product"`
	MrID    int    `json:"mr_id"`
	Author  string `json:"author"`
	//Merger  string `json:"merger"` // do we need merged by whom info?
	State string `json:"state"`
}

type JenkinsBuildExporter struct {
	User        string            `json:"user"`
	BuildUrl    string            `json:"build_url"`
	BuildCause  string            `json:"build_cause"`
	JobName     string            `json:"job_name"`
	BuildNum    string            `json:"build_num"`
	RepoUrl     string            `json:"repo_url"`
	Branch      string            `json:"branch"`
	Sha1        string            `json:"SHA1"`
	Timestamp   int64             `json:"timestamp"`
	Result      string            `json:"result"`
	Bugs        map[string]string `json:"bugs"`
	IsTestJob   bool              `json:"is_test_job"`
	TestDetails TestDetails       `json:"test_details"`
}

func (be *JenkinsBuildExporter) LinkBug(bug string) {
	if be.Bugs == nil {
		be.Bugs = map[string]string{}
	}
	be.Bugs[bug] = ""
}

type TestDetails struct {
	TestedRepo   string `json:"tested_repo"`
	TestedBranch string `json:"tested_branch"`
	TestedSha1   string `json:"tested_commit"`
	TotalCount   int    `json:"total_count"`
	FailCount    int    `json:"fail_count"`
	SkipCount    int    `json:"skip_count"`
}

type JenkinsBuildDetails struct {
	Class   string `json:"_class"`
	Actions []struct {
		Class      string `json:"_class,omitempty"`
		Parameters []struct {
			Class string `json:"_class"`
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"parameters,omitempty"`
		Causes []struct {
			Class            string `json:"_class"`
			ShortDescription string `json:"shortDescription"`
			UserID           string `json:"userId"`
			UserName         string `json:"userName"`
		} `json:"causes,omitempty"`
		LastBuiltRevision struct {
			Sha1   string `json:"SHA1"`
			Branch []struct {
				Sha1 string `json:"SHA1"`
				Name string `json:"name"`
			} `json:"branch"`
		} `json:"lastBuiltRevision,omitempty"`
		RemoteUrls []string `json:"remoteUrls,omitempty"`
		ScmName    string   `json:"scmName,omitempty"`
		FailCount  int      `json:"failCount,omitempty"`
		SkipCount  int      `json:"skipCount,omitempty"`
		TotalCount int      `json:"totalCount,omitempty"`
		URLName    string   `json:"urlName,omitempty"`
	} `json:"actions"`
	Artifacts []struct {
		DisplayPath  string `json:"displayPath"`
		FileName     string `json:"fileName"`
		RelativePath string `json:"relativePath"`
	} `json:"artifacts"`
	Building          bool        `json:"building"`
	Description       interface{} `json:"description"`
	DisplayName       string      `json:"displayName"`
	Duration          int         `json:"duration"`
	EstimatedDuration int         `json:"estimatedDuration"`
	Executor          interface{} `json:"executor"`
	FullDisplayName   string      `json:"fullDisplayName"`
	ID                string      `json:"id"`
	KeepLog           bool        `json:"keepLog"`
	Number            int         `json:"number"`
	QueueID           int         `json:"queueId"`
	Result            string      `json:"result"`
	Timestamp         int64       `json:"timestamp"`
	URL               string      `json:"url"`
	BuiltOn           string      `json:"builtOn"`
	ChangeSet         struct {
		Class string `json:"_class"`
		Items []struct {
			Class         string   `json:"_class"`
			AffectedPaths []string `json:"affectedPaths"`
			Author        struct {
				AbsoluteURL string `json:"absoluteUrl"`
				FullName    string `json:"fullName"`
			} `json:"author"`
			AuthorEmail string `json:"authorEmail"`
			Comment     string `json:"comment"`
			CommitID    string `json:"commitId"`
			Date        string `json:"date"`
			ID          string `json:"id"`
			Msg         string `json:"msg"`
			Paths       []struct {
				EditType string `json:"editType"`
				File     string `json:"file"`
			} `json:"paths"`
			Timestamp int64 `json:"timestamp"`
		} `json:"items"`
		Kind string `json:"kind"`
	} `json:"changeSet"`
	Culprits []interface{} `json:"culprits"`
}
