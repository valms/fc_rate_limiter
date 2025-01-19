package ratelimiter

import (
	"fmt"
	"github.com/valms/fc_rate_limiter/internal/config"
	"github.com/valms/fc_rate_limiter/internal/repository"
)

type RateLimiter struct {
	repository    repository.RateLimitRepository
	ipLimit       int
	apiKeyLimit   int
	blockDuration int
}

func NewLimiterService(repo repository.RateLimitRepository, loadConfig *config.Config) (*RateLimiter, error) {
	return &RateLimiter{
		repository:    repo,
		ipLimit:       loadConfig.RateLimit.IPLimit,
		apiKeyLimit:   loadConfig.RateLimit.APIKeyLimit,
		blockDuration: loadConfig.RateLimit.BlockDuration,
	}, nil
}

func (l *RateLimiter) IsRateLimitByIP(ip string) (bool, error) {
	blocked, err := l.repository.KeyExists(l.blockedKey(ip))
	fmt.Printf("[Rate Limiter] IP %s - Block Status: %v\n", ip, blocked)

	if err != nil {
		return true, fmt.Errorf("error checking block status: %w", err)
	}

	if blocked {
		fmt.Printf("[Rate Limiter] IP %s is BLOCKED for next %d seconds\n",
			ip, l.blockDuration)
		return true, nil
	}

	currentCount, err := l.repository.GetCounter(l.requestKey(ip))
	if err != nil {
		return false, fmt.Errorf("error getting current count: %w", err)
	}

	fmt.Printf("[Rate Limiter] IP %s - Current Usage: %d/%d\n",
		ip, currentCount, l.ipLimit)

	if currentCount >= l.ipLimit {
		fmt.Printf("[Rate Limiter] IP %s exceeded limit. Blocking for %d seconds\n",
			ip, l.blockDuration)

		err = l.repository.SetExpiration(l.blockedKey(ip), l.blockDuration)
		if err != nil {
			return false, fmt.Errorf("error setting block: %w", err)
		}
		return true, nil
	}

	newCount, err := l.repository.IncrementCounter(l.requestKey(ip))
	if err != nil {
		return false, fmt.Errorf("error incrementing counter: %w", err)
	}

	fmt.Printf("[Rate Limiter] IP %s - Updated Usage: %d/%d\n",
		ip, newCount, l.ipLimit)

	if newCount == 1 {
		fmt.Printf("[Rate Limiter] IP %s - First request, setting window expiration\n", ip)
		if err := l.repository.SetExpiration(l.requestKey(ip), 20); err != nil {
			return false, fmt.Errorf("error setting window: %w", err)
		}
	}

	if newCount >= l.ipLimit {
		fmt.Printf("[Rate Limiter] IP %s reached limit. Blocking for %d seconds\n",
			ip, l.blockDuration)

		if err := l.repository.SetExpiration(l.blockedKey(ip), l.blockDuration); err != nil {
			return false, fmt.Errorf("error setting block: %w", err)
		}
		return true, nil
	}

	return false, nil
}

func (l *RateLimiter) IsRateLimitByToken(token string) (bool, error) {
	blocked, err := l.repository.KeyExists(l.blockedKey(token))
	fmt.Printf("[Rate Limiter] API_KEY %s - Block Status: %v\n", token, blocked)

	if err != nil {
		return true, fmt.Errorf("error checking block status: %w", err)
	}

	if blocked {
		fmt.Printf("[Rate Limiter] API_KEY %s is BLOCKED for next %d seconds\n",
			token, l.blockDuration)
		return true, nil
	}

	currentCount, err := l.repository.GetCounter(l.requestKey(token))
	if err != nil {
		return false, fmt.Errorf("error getting current count: %w", err)
	}

	fmt.Printf("[Rate Limiter] API_KEY %s - Current Usage: %d/%d\n",
		token, currentCount, l.apiKeyLimit)

	if currentCount >= l.apiKeyLimit {
		fmt.Printf("[Rate Limiter] API_KEY %s exceeded limit. Blocking for %d seconds\n",
			token, l.blockDuration)

		err = l.repository.SetExpiration(l.blockedKey(token), l.blockDuration)
		if err != nil {
			return false, fmt.Errorf("error setting block: %w", err)
		}
		return true, nil
	}

	newCount, err := l.repository.IncrementCounter(l.requestKey(token))
	if err != nil {
		return false, fmt.Errorf("error incrementing counter: %w", err)
	}

	fmt.Printf("[Rate Limiter] API_KEY %s - Updated Usage: %d/%d\n",
		token, newCount, l.apiKeyLimit)

	if newCount == 1 {
		fmt.Printf("[Rate Limiter] API_KEY %s - First request, setting window expiration\n", token)
		if err := l.repository.SetExpiration(l.requestKey(token), 20); err != nil {
			return false, fmt.Errorf("error setting window: %w", err)
		}
	}

	if newCount >= l.apiKeyLimit {
		fmt.Printf("[Rate Limiter] API_KEY %s reached limit. Blocking for %d seconds\n",
			token, l.blockDuration)

		if err := l.repository.SetExpiration(l.blockedKey(token), l.blockDuration); err != nil {
			return false, fmt.Errorf("error setting block: %w", err)
		}
		return true, nil
	}

	return false, nil
}

func (l *RateLimiter) requestKey(identifier string) string {
	return fmt.Sprintf("req_count:%s", identifier)
}

func (l *RateLimiter) blockedKey(identifier string) string {
	return fmt.Sprintf("block:%s", identifier)
}
