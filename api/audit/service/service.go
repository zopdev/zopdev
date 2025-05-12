package service

import (
	"time"

	"gofr.dev/pkg/gofr"
	gofrHttp "gofr.dev/pkg/gofr/http"

	"github.com/zopdev/zopdev/api/audit/client"
	"github.com/zopdev/zopdev/api/audit/rules/overprovision"
	"github.com/zopdev/zopdev/api/audit/store"
)

// Service is a struct that holds the rules and their execution logic.
// It is responsible for executing the rules and returning the results.
type Service struct {
	rules           map[string]Rule
	categoryRuleMap map[string][]Rule

	store Store
}

func New(str Store) *Service {
	s := &Service{
		store: str,

		rules:           make(map[string]Rule),
		categoryRuleMap: make(map[string][]Rule),
	}

	// Register rules here
	s.rules["sql_instance_peak"] = &overprovision.SQLInstancePeak{}

	// parse the added rules and create a map of category to rules
	// This is done to avoid parsing the rule map to avoid iterating over all rules
	// every time we need to execute category specific rules
	s.parse()

	return s
}

func (s *Service) parse() {
	for _, rule := range s.rules {
		category := rule.GetCategory()
		if _, exists := s.categoryRuleMap[category]; !exists {
			s.categoryRuleMap[category] = make([]Rule, 0)
		}

		s.categoryRuleMap[category] = append(s.categoryRuleMap[category], rule)
	}
}

// RunByID executes the rule with the given ruleID and cloudAccId. It fetches the cloud credentials from the cloud-account entity
// and passes it to the rule for execution.
func (s *Service) RunByID(ctx *gofr.Context, ruleID string, cloudAccID int64) (*store.Result, error) {
	rule, exists := s.rules[ruleID]
	if !exists {
		return nil, gofrHttp.ErrorEntityNotFound{Name: "Rule", Value: ruleID}
	}

	ca, err := client.GetCloudCredentials(ctx, cloudAccID)
	if err != nil {
		return nil, err
	}

	// create a result entry in the database
	res, err := s.store.CreatePending(ctx, &store.Result{
		RuleID:         ruleID,
		CloudAccountID: cloudAccID,
		Result:         &store.ResultData{},
		EvaluatedAt:    time.Now(),
	})
	if err != nil {
		return nil, err
	}

	result, err := rule.Execute(ctx, ca)
	if err != nil {
		return nil, err
	}

	// update the result entry in the database
	res.Result.Data = result

	err = s.store.UpdateResult(ctx, res)
	if err != nil {
		return res, err
	}

	return res, nil
}

// RunByCategory executes all the rules in the given category and returns the results.
func (s *Service) RunByCategory(ctx *gofr.Context, category string, cloudAccID int64) ([]*store.Result, error) {
	results := make([]*store.Result, 0)

	ca, err := client.GetCloudCredentials(ctx, cloudAccID)
	if err != nil {
		return nil, err
	}

	rules, exists := s.categoryRuleMap[category]
	if !exists {
		return nil, gofrHttp.ErrorEntityNotFound{Name: "Category", Value: category}
	}

	for _, rule := range rules {
		// create a result entry in the database
		res, er := s.store.CreatePending(ctx, &store.Result{
			RuleID:         rule.GetName(),
			CloudAccountID: cloudAccID,
			Result:         &store.ResultData{},
			EvaluatedAt:    time.Now(),
		})
		if er != nil {
			ctx.Errorf("error creating result entry: %v", er)
			continue
		}

		result, er := rule.Execute(ctx, ca)
		if er != nil {
			return nil, err
		}

		// update the result entry in the database
		res.Result.Data = result
		_ = s.store.UpdateResult(ctx, res)

		results = append(results, res)
	}

	return results, nil
}

// RunAll executes all the rules in the rule engine and returns the results.
// It fetches the cloud credentials from the cloud-account entity and passes it to each rule for execution.
// It returns a slice of ResultData, which contains the results of each rule executed.
// The results are grouped by category.
func (s *Service) RunAll(ctx *gofr.Context, cloudAccID int64) (map[string][]*store.Result, error) {
	results := make(map[string][]*store.Result)

	ca, err := client.GetCloudCredentials(ctx, cloudAccID)
	if err != nil {
		return nil, err
	}

	for _, rule := range s.rules {
		// create a result entry in the database
		res, er := s.store.CreatePending(ctx, &store.Result{
			RuleID:         rule.GetName(),
			CloudAccountID: cloudAccID,
			Result:         &store.ResultData{},
			EvaluatedAt:    time.Now(),
		})
		if er != nil {
			ctx.Errorf("error creating result entry: %v", er)
			continue
		}

		result, er := rule.Execute(ctx, ca)
		if er != nil {
			return nil, er
		}

		// update the result entry in the database
		res.Result.Data = result
		_ = s.store.UpdateResult(ctx, res)

		_, ok := results[rule.GetCategory()]
		if !ok {
			results[rule.GetCategory()] = make([]*store.Result, 0)
		}

		results[rule.GetCategory()] = append(results[rule.GetCategory()], res)
	}

	return results, nil
}

// GetResultByID retrieves the result of a specific rule execution by its ID.
func (s *Service) GetResultByID(ctx *gofr.Context, cloudAccID int64, ruleID string) (*store.Result, error) {
	res, err := s.store.GetLastRun(ctx, cloudAccID, ruleID)
	if err != nil {
		return nil, err
	}

	if res == nil {
		return nil, gofrHttp.ErrorEntityNotFound{Name: "Result", Value: ruleID}
	}

	return res, nil
}

// GetAllResults retrieves the latest results for all rules associated with a given cloud account ID.
// It organizes the results into a map where the keys are rule categories and the values are slices
// of results belonging to those categories.
func (s *Service) GetAllResults(ctx *gofr.Context, cloudAccID int64) (map[string][]*store.Result, error) {
	result := make(map[string][]*store.Result)

	for _, rule := range s.rules {
		res, err := s.store.GetLastRun(ctx, cloudAccID, rule.GetName())
		if err != nil {
			return nil, err
		}

		if res == nil {
			continue
		}

		if _, ok := result[rule.GetCategory()]; !ok {
			result[rule.GetCategory()] = make([]*store.Result, 0)
		}

		result[rule.GetCategory()] = append(result[rule.GetCategory()], res)
	}

	return result, nil
}
