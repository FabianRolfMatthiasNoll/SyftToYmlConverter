package provider

import (
	"context"
	"encoding/base64"
	"regexp"
	"time"

	"syfttoymlconverter/internal/model"

	"github.com/google/go-github/v37/github"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
)

var (
	//nolint:gochecknoglobals // singleton instance for github
	githubClient = newGithub("33e1d1373efba32ce669ff782812f6aed422714c")

	// regex to split github address into owner/reponame
	githubRegEx = regexp.MustCompile(`^github\.com/([^/]+)/([^/]+)$`)
)

type githubProvider struct {
	client *github.Client
}

func newGithub(token string) *githubProvider {
	oauthClient := oauth2.NewClient(
		context.Background(),
		oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		),
	)

	return &githubProvider{
		client: github.NewClient(oauthClient),
	}
}

func (g *githubProvider) getInfo(path, tag string) model.RepoInfo {
	matches := githubRegEx.FindStringSubmatch(path)
	if matches == nil {
		return model.RepoInfo{}
	}

	owner := matches[1]
	reponame := matches[2]

	info := model.RepoInfo{}

	repo, ok := g.getRepoData(owner, reponame)
	if ok {
		info.FullName = repo.GetFullName()
		info.Description = repo.GetDescription()

		// github returns NOASSERTION when license unknown
		if repo.License != nil && repo.License.GetSPDXID() != "NOASSERTION" {
			info.SPDX = repo.License.GetSPDXID()
		}
	}

	releaseDate, ok := g.getReleaseDate(owner, reponame, tag)
	if ok {
		info.Release = releaseDate
	}

	return info
}

func (g *githubProvider) getRepoData(owner, reponame string) (*github.Repository, bool) {
	repo, _, err := g.client.Repositories.Get(context.Background(), owner, reponame)
	if err != nil {
		log.Warn().Err(err).Msgf("Failed to fetch repository data for %s/%s", owner, reponame)

		return nil, false
	}

	return repo, true
}

func (g *githubProvider) getReleaseDate(owner, reponame, tag string) (time.Time, bool) {
	release, ok := g.getReleaseDateByRelease(owner, reponame, tag)
	if ok {
		return release, ok
	}

	release, ok = g.getReleaseDateByTag(owner, reponame, tag)
	if ok {
		return release, ok
	}

	return time.Time{}, false
}

func (g *githubProvider) getReleaseDateByRelease(owner, reponame, tag string) (time.Time, bool) {
	repo, _, err := g.client.Repositories.GetReleaseByTag(context.Background(), owner, reponame, tag)
	if err != nil || repo.PublishedAt == nil {
		log.Debug().Err(err).Msgf("Failed to fetch release date by tag for %s/%s@%s", owner, reponame, tag)

		return time.Time{}, false
	}

	return repo.GetPublishedAt().Local(), true
}

func (g *githubProvider) getReleaseDateByTag(owner, reponame, tag string) (time.Time, bool) {
	repoTags, _, err := g.client.Repositories.ListTags(context.Background(), owner, reponame, nil)
	if err != nil {
		log.Debug().Err(err).Msgf("Failed to fetch tags for %s/%s", owner, reponame)

		return time.Time{}, false
	}

	for _, repoTag := range repoTags {
		if repoTag.GetName() == tag {
			sha := repoTag.GetCommit().GetSHA()

			commit, _, err := g.client.Repositories.GetCommit(context.Background(), owner, reponame, sha)
			if err != nil {
				log.Debug().Err(err).Msgf("Failed to fetch commit info for %s/%s@%s", owner, reponame, sha)

				return time.Time{}, false
			}

			return commit.GetCommit().GetAuthor().GetDate(), true
		}
	}

	log.Debug().Msgf("Failed to fetch release date for %s/%s@%s", owner, reponame, tag)

	return time.Time{}, false
}

func (g *githubProvider) getLicenseFromRepo(source string) (string, error) {
	matches := githubRegEx.FindStringSubmatch(source)
	if matches == nil {
		return "", errors.Errorf("pattern missmatch, no github source [%s]", source)
	}

	owner := matches[1]
	reponame := matches[2]

	repoLicense, _, err := g.client.Repositories.License(context.Background(), owner, reponame)
	if err != nil {
		return "", errors.Wrap(err, "failed to fetch license from github")
	}

	if repoLicense.GetEncoding() != "base64" {
		return "", errors.Errorf("unexpected license encoding [%s]", repoLicense.GetEncoding())
	}

	str, err := base64.StdEncoding.DecodeString(repoLicense.GetContent())
	if err != nil {
		return "", errors.Wrap(err, "failed to decode base64 license text")
	}

	return string(str), nil
}

func (g *githubProvider) getSpdxLicense(spdx string) (string, error) {
	license, _, err := g.client.Licenses.Get(context.Background(), spdx)
	if err != nil {
		return "", errors.Wrap(err, "fao;ed to fetch spdx data from github")
	}

	return license.GetBody(), nil
}
