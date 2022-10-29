package git_utils

import (
	"errors"
	"fmt"
	"github.com/luohuahuang/qex/config"
	"github.com/luohuahuang/qex/pkg/mattermost"
	"net/http"
	"time"

	"github.com/xanzy/go-gitlab"
)

func QueryGitlabProjectMRs(client *gitlab.Client, projectId int, startTime *time.Time, endTime *time.Time) ([]*gitlab.MergeRequest, error) {
	pageNum := 1
	totalMRs := make([]*gitlab.MergeRequest, 0)

	for {
		opt := &gitlab.ListProjectMergeRequestsOptions{
			Scope: gitlab.String("all"),
		}

		if startTime != nil && endTime != nil {
			opt.CreatedAfter = gitlab.Time(*startTime)
			opt.CreatedBefore = gitlab.Time(*endTime)
		}

		opt.PerPage = 20
		opt.Page = pageNum

		mrs, resp, err := client.MergeRequests.ListProjectMergeRequests(projectId, opt)

		if resp != nil && (resp.StatusCode == http.StatusNotFound ||
			resp.StatusCode == http.StatusUnauthorized ||
			resp.StatusCode == http.StatusForbidden) {
			mrs = []*gitlab.MergeRequest{}
			err = nil
		}

		if err != nil {
			mattermost.SendAlert(err, config.MatterMostMonitor)
			return nil, err
		}

		totalMRs = append(totalMRs, mrs...)
		pageNum = pageNum + 1

		if len(mrs) <= 0 {
			break
		}
	}

	return totalMRs, nil
}

func QueryGitlabMRByBranches(client *gitlab.Client, projectId int, iid int) (*gitlab.MergeRequest, error) {

	opt := &gitlab.ListProjectMergeRequestsOptions{
		IIDs: &[]int{iid},
	}

	mrs, resp, err := client.MergeRequests.ListProjectMergeRequests(projectId, opt)

	if resp != nil && (resp.StatusCode == http.StatusNotFound ||
		resp.StatusCode == http.StatusUnauthorized ||
		resp.StatusCode == http.StatusForbidden) {
		mattermost.SendAlert(errors.New(fmt.Sprintf("fail to get MR info for project: %d, %v", projectId, opt)), config.MatterMostMonitor)
		return nil, err
	}

	if err != nil {
		mattermost.SendAlert(err, config.MatterMostMonitor)
		return nil, err
	}

	if len(mrs) > 0 {
		return mrs[0], nil
	}
	return nil, err
}

func QueryGitlabProjectEvents(client *gitlab.Client, gitLabId uint32, startTime time.Time, endTime time.Time) ([]*gitlab.ContributionEvent, error) {
	pageNum := 1
	totalEvents := make([]*gitlab.ContributionEvent, 0)

	newStartTime := startTime.Add(-24 * time.Hour)

	for {
		after := gitlab.ISOTime(newStartTime)
		before := gitlab.ISOTime(endTime)
		sort := "asc"

		opt := &gitlab.ListContributionEventsOptions{
			After:  &after,
			Before: &before,
			Sort:   &sort,
		}

		opt.PerPage = 20
		opt.Page = pageNum

		events, resp, err := client.Events.ListProjectVisibleEvents(int(gitLabId), opt)

		if resp != nil && (resp.StatusCode == http.StatusNotFound ||
			resp.StatusCode == http.StatusUnauthorized ||
			resp.StatusCode == http.StatusForbidden) {
			events = []*gitlab.ContributionEvent{}
			err = nil
		}

		if err != nil {
			mattermost.SendAlert(err, config.MatterMostMonitor)
			return nil, err
		}

		totalEvents = append(totalEvents, events...)
		pageNum = pageNum + 1

		if len(events) <= 0 {
			break
		}
	}

	return totalEvents, nil
}
