package repository

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"task4.2.3/internal/models"
	"task4.2.3/internal/monitoring"
)

type AddressRepositoryProxy struct {
	repo  AddressRepository
	cache *redis.Client
	ttl   time.Duration
}

func NewAddressRepositoryProxy(repo AddressRepository, cache *redis.Client, ttl time.Duration) AddressRepository {
	return &AddressRepositoryProxy{
		repo:  repo,
		cache: cache,
		ttl:   ttl,
	}
}

func (p *AddressRepositoryProxy) Search(query string) ([]models.Address, error) {
	ctx := context.Background()
	key := "search:" + query

	cacheStartTime := time.Now()
	cached, err := p.cache.Get(ctx, key).Result()
	cacheDuration := time.Since(cacheStartTime).Seconds()

	monitoring.CacheRequestDuration.WithLabelValues("Search").Observe(cacheDuration)

	if err == nil {
		log.Println("Data is in cache!")
		var addresses []models.Address
		if json.Unmarshal([]byte(cached), &addresses) == nil {
			return addresses, nil
		}
	}
	log.Println("Data not in cache!")

	apiStartTime := time.Now()
	data, err := p.repo.Search(query)
	apiDuration := time.Since(apiStartTime).Seconds()

	monitoring.ExternalAPIRequestDuration.WithLabelValues("Search").Observe(apiDuration)

	if err != nil {
		log.Println("Error querying data in repo:", err)
		return nil, err
	}

	if jsonStr, err := json.Marshal(data); err == nil {
		if err := p.cache.Set(ctx, key, jsonStr, p.ttl).Err(); err != nil {
			log.Println("Error caching search result:", err)
		}
	}

	return data, nil
}

func (p *AddressRepositoryProxy) Geocode(lat, lng string) ([]models.Address, error) {
	ctx := context.Background()
	key := "geocode:" + lat + ":" + lng

	cacheStartTime := time.Now()
	cached, err := p.cache.Get(ctx, key).Result()
	cacheDuration := time.Since(cacheStartTime).Seconds()

	monitoring.CacheRequestDuration.WithLabelValues("Geocode").Observe(cacheDuration)

	if err == nil {
		log.Println("Data is in cache!")
		var addresses []models.Address
		if json.Unmarshal([]byte(cached), &addresses) == nil {
			return addresses, nil
		}
	}
	log.Println("Data not in cache!")

	apiStartTime := time.Now()
	data, err := p.repo.Geocode(lat, lng)
	apiDuration := time.Since(apiStartTime).Seconds()

	monitoring.ExternalAPIRequestDuration.WithLabelValues("Geocode").Observe(apiDuration)

	if err != nil {
		log.Println("Error querying data in repo:", err)
		return nil, err
	}

	if jsonStr, err := json.Marshal(data); err == nil {
		if err := p.cache.Set(ctx, key, jsonStr, p.ttl).Err(); err != nil {
			log.Println("Error caching geocode result:", err)
		}
	}

	return data, nil
}
