package extraction

import (
	"context"
	"deprec/logging"
	"errors"
	"fmt"
	"github.com/google/go-github/v48/github"
	"github.com/thoas/go-funk"
	"reflect"
	"strings"
)

func BatchQuery[T any](ctx context.Context, client *GitHubClientWrapper, queries map[string]string, vars map[string]any) (map[string]T, error) {
	// Create a query using reflection (see https://github.com/shurcooL/githubv4/issues/17)
	// for when we don't know the exact query before runtime.
	var t T
	fieldType := reflect.TypeOf(t)
	var fields []reflect.StructField
	fieldToKey := map[string]string{}
	idx := 0
	for key, q := range queries {
		name := fmt.Sprintf("Field%d", idx)
		fieldToKey[name] = key
		fields = append(fields, reflect.StructField{
			Name: name,
			Type: fieldType,
			Tag:  reflect.StructTag(fmt.Sprintf(`graphql:"field%d:%s"`, idx, strings.ReplaceAll(q, "\"", "\\\""))),
		})
		idx++
	}
	// TODO: an upper bound should be added
	if len(fields) == 0 {
		return nil, errors.New("no query to run")
	}
	q := reflect.New(reflect.StructOf(fields)).Elem()
	if err := client.Client.GraphQL().Query(ctx, q.Addr().Interface(), vars); err != nil {
		return nil, err
	}
	res := map[string]T{}
	for _, sf := range reflect.VisibleFields(q.Type()) {
		key := fieldToKey[sf.Name]
		v := q.FieldByIndex(sf.Index)
		res[key] = v.Interface().(T)
	}
	return res, nil
}

func FetchContributorInfo(ctx context.Context, contributors []*github.Contributor, c *GitHubClientWrapper) (map[string]ContributorInfo, error) {

	logging.Logger.Info("HITTING V4 API")

	// Doing this over REST would take O(n) requests, using GraphQL takes O(1).
	userQueries := map[string]string{}
	for i, contributor := range contributors {
		login := contributor.GetLogin()
		if login == "" {
			continue
		}
		if strings.HasSuffix(login, "[bot]") {
			continue
		}
		userQueries[fmt.Sprint(i)] = fmt.Sprintf("user(login:\"%s\")", login)
	}
	if len(userQueries) == 0 {
		// We didn't add any users.
		return nil, errors.New("no contributors to fetch info for")
	}
	r, err := BatchQuery[ContributorInfo](ctx, c, userQueries, map[string]any{})
	if err != nil {
		return nil, err
	}

	m := funk.Map(r, func(s string, q ContributorInfo) (string, ContributorInfo) { return q.Login, q }).(map[string]ContributorInfo)
	return m, nil
}

type ContributorInfo struct {
	Repositories struct {
		TotalCount int
	}
	Sponsors struct {
		TotalCount int
	}
	Organizations struct {
		TotalCount int
	}
	Company string
	Login   string
}

func errorTooManyContributors(err error) bool {
	if err == nil {
		return false
	}
	var e *github.ErrorResponse
	ok := errors.As(err, &e)
	if !ok {
		return false
	}
	return e.Response.StatusCode == 403 && strings.Contains(e.Message, "list is too large")
}
