package github

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/jenkins-x/go-scm/scm"
	"github.com/pkg/errors"
)

// ListIssueEvents list issue events
func (c *Client) ListIssueEvents(org, repo string, number int) ([]*scm.ListedIssueEvent, error) {
	ctx := context.Background()
	fullName := c.repositoryName(org, repo)
	events, _, err := c.client.Issues.ListEvents(ctx, fullName, number, c.createListOptions())
	return events, err
}

// AssignIssue assigns issue
func (c *Client) AssignIssue(owner, repo string, number int, logins []string) error {
	ctx := context.Background()
	fullName := c.repositoryName(owner, repo)
	_, err := c.client.Issues.AssignIssue(ctx, fullName, number, logins)
	return err
}

// UnassignIssue unassigns issue
func (c *Client) UnassignIssue(owner, repo string, number int, logins []string) error {
	ctx := context.Background()
	fullName := c.repositoryName(owner, repo)
	_, err := c.client.Issues.UnassignIssue(ctx, fullName, number, logins)
	return err
}

// AddLabel adds a label
func (c *Client) AddLabel(owner, repo string, number int, label string) error {
	ctx := context.Background()
	fullName := c.repositoryName(owner, repo)
	_, err := c.client.Issues.AddLabel(ctx, fullName, number, label)
	return err
}

// RemoveLabel removes labesl
func (c *Client) RemoveLabel(owner, repo string, number int, label string) error {
	ctx := context.Background()
	fullName := c.repositoryName(owner, repo)
	_, err := c.client.Issues.DeleteLabel(ctx, fullName, number, label)
	return err
}

// DeleteComment delete comments
func (c *Client) DeleteComment(org, repo string, number, ID int) error {
	ctx := context.Background()
	fullName := c.repositoryName(org, repo)
	_, err := c.client.Issues.DeleteComment(ctx, fullName, number, ID)
	return err
}

// DeleteStaleComments iterates over comments on an issue/PR, deleting those which the 'isStale'
// function identifies as stale. If 'comments' is nil, the comments will be fetched from GitHub.
func (c *Client) DeleteStaleComments(org, repo string, number int, comments []*scm.Comment, isStale func(*scm.Comment) bool) error {
	var err error
	if comments == nil {
		comments, err = c.ListIssueComments(org, repo, number)
		if err != nil {
			return fmt.Errorf("failed to list comments while deleting stale comments. err: %v", err)
		}
	}
	for _, comment := range comments {
		if isStale(comment) {
			if err := c.DeleteComment(org, repo, number, comment.ID); err != nil {
				return fmt.Errorf("failed to delete stale comment with ID '%d'", comment.ID)
			}
		}
	}
	return nil
}

// ListIssueComments list comments associated with an issue
func (c *Client) ListIssueComments(org, repo string, number int) ([]*scm.Comment, error) {
	ctx := context.Background()
	fullName := c.repositoryName(org, repo)
	comments, _, err := c.client.Issues.ListComments(ctx, fullName, number, c.createListOptions())
	return comments, err
}

// GetIssueLabels returns the issue labels
func (c *Client) GetIssueLabels(org, repo string, number int) ([]*scm.Label, error) {
	ctx := context.Background()
	fullName := c.repositoryName(org, repo)
	labels, _, err := c.client.Issues.ListLabels(ctx, fullName, number, c.createListOptions())
	return labels, err
}

// CreateComment create a comment
func (c *Client) CreateComment(owner, repo string, number int, pr bool, comment string) error {
	fullName := c.repositoryName(owner, repo)
	commentInput := scm.CommentInput{
		Body: comment,
	}
	ctx := context.Background()
	if pr {
		_, response, err := c.client.PullRequests.CreateComment(ctx, fullName, number, &commentInput)
		if err != nil {
			var b bytes.Buffer
			_, cperr := io.Copy(&b, response.Body)
			if cperr != nil {
				return errors.Wrapf(cperr, "response: %s", b.String())
			}
			return errors.Wrapf(err, "response: %s", b.String())
		}

	} else {
		_, response, err := c.client.Issues.CreateComment(ctx, fullName, number, &commentInput)
		if err != nil {
			var b bytes.Buffer
			_, cperr := io.Copy(&b, response.Body)
			if cperr != nil {
				return errors.Wrapf(cperr, "reponse: %s", b.String())
			}
			return errors.Wrapf(err, "response: %s", b.String())
		}
	}
	return nil
}

// ReopenIssue reopen an issue
func (c *Client) ReopenIssue(owner, repo string, number int) error {
	panic("implement me")
}

// FindIssues find issues
func (c *Client) FindIssues(query, sort string, asc bool) ([]scm.Issue, error) {
	panic("implement me")
}

// CloseIssue close issue
func (c *Client) CloseIssue(owner, repo string, number int) error {
	panic("implement me")
}
