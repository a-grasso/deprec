package githubapi

import (
	"context"
	"deprec/cache"
	"deprec/model"
	"errors"
	"fmt"
	"github.com/google/go-github/v48/github"
	"github.com/thoas/go-funk"
	"reflect"
	"strings"
)

func BatchQuery[T any](ctx context.Context, client *Client, queries map[string]string, vars map[string]any) (map[string]T, error) {
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
	if err := client.GraphQL().Query(ctx, q.Addr().Interface(), vars); err != nil {
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

func (ql *GraphQLWrapper) FetchContributorInfo(ctx context.Context, repo string, contributors []*github.Contributor) (map[string]model.ContributorInfo, error) {

	coll := ql.Cache.Database("query_contributor_info").Collection(repo)

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
		return nil, errors.New("no contributors to fetch info for")
	}

	batchQuery := func() (map[string]model.ContributorInfo, error) {
		return BatchQuery[model.ContributorInfo](ctx, ql.Client, userQueries, map[string]any{})
	}

	info, err := cache.FetchBatchQuery[model.ContributorInfo](ctx, coll, batchQuery)

	mapped := funk.Map(info, func(q model.ContributorInfo) (string, model.ContributorInfo) { return q.Login, q }).(map[string]model.ContributorInfo)

	return mapped, err
}
