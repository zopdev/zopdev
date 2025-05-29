package service

import (
	"sync"

	"gofr.dev/pkg/gofr"

	"github.com/zopdev/zopdev/api/resources/client"
)

func (s *Service) SyncCron(ctx *gofr.Context) {
	cl, err := s.http.GetAllCloudAccounts(ctx)
	if err != nil {
		ctx.Metrics().IncrementCounter(ctx, "sync_error_count")
		ctx.Errorf("failed to get cloud accounts: %v", err)

		return
	}

	errs := s.SyncAll(ctx, cl)
	if len(errs) > 0 {
		ctx.Metrics().IncrementCounter(ctx, "sync_error_count")
	}
}

func (s *Service) SyncAll(ctx *gofr.Context, accounts []client.CloudAccount) []error {
	var (
		wg  sync.WaitGroup
		err []error
	)

	for _, account := range accounts {
		wg.Add(1)

		go func(accountID int64) {
			defer wg.Done()

			_, er := s.SyncResources(ctx, accountID)
			if er != nil {
				ctx.Errorf("failed to sync resources for account %d: %v", accountID, er)
				err = append(err, er)
			}
		}(account.ID)
	}

	wg.Wait()

	return err
}
