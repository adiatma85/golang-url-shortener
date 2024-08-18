package url

import (
	"context"
	"fmt"
	"strconv"

	"github.com/adiatma85/golang-url-shortener/src/business/entity"
	"github.com/adiatma85/own-go-sdk/codes"
	"github.com/adiatma85/own-go-sdk/errors"
	"github.com/adiatma85/own-go-sdk/query"
)

// Implement this
func (u *url) AssignCounterScheduler(ctx context.Context) error {
	// Pertama adalah baca semua keys yang ada di Redis
	scanKey := fmt.Sprintf(entity.UrlCountingRedisKey, "*")
	allKeys, err := u.redis.Scan(ctx, scanKey)
	if err != nil {
		return err
	}

	for _, singleKey := range allKeys {
		count, err := u.redis.Get(ctx, singleKey)

		if err != nil {
			u.log.Error(ctx, errors.NewWithCode(codes.CodeRedisGet, fmt.Sprintf("error when processing redis with key %s", singleKey)))
			continue
		}

		num, err := strconv.Atoi(count)
		if err != nil {
			u.log.Error(ctx, err)
			continue
		}

		urlParam := entity.UrlParam{
			ShortenUrl: singleKey,
			QueryOption: query.Option{
				IsActive: true,
			},
		}

		url, err := u.url.Get(ctx, urlParam)
		if err != nil {
			u.log.Error(ctx, err)
			continue
		}

		updateParam := entity.UpdateUrlParam{
			Visit: url.Visit + int64(num),
		}

		// Update to the database
		if err = u.url.Update(ctx, updateParam, urlParam); err != nil {
			return err
		}

		// Decrement the Redis
		if err = u.redis.DecrementBy(ctx, singleKey, int64(num)); err != nil {
			u.log.Error(ctx, err)
			continue
		}
	}

	return nil
}
