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
	fmt.Printf("[Rate Limiter] IP %s - Status de Bloqueio: %v\n", ip, blocked)

	if err != nil {
		return true, fmt.Errorf("erro ao verificar status de bloqueio: %w", err)
	}

	if blocked {
		fmt.Printf("[Rate Limiter] IP %s está bloqueado pelos próximos %d segundos\n",
			ip, l.blockDuration)
		return true, nil
	}

	currentCount, err := l.repository.GetCounter(l.requestKey(ip))
	if err != nil {
		return false, fmt.Errorf("erro ao obter contagem atual: %w", err)
	}

	fmt.Printf("[Rate Limiter] IP %s - Uso Atual: %d/%d\n",
		ip, currentCount, l.ipLimit)

	if currentCount > l.ipLimit {
		fmt.Printf("[Rate Limiter] IP %s excedeu o limite. Bloqueando por %d segundos\n",
			ip, l.blockDuration)

		err = l.repository.SetExpiration(l.blockedKey(ip), l.blockDuration)
		if err != nil {
			return false, fmt.Errorf("erro ao definir bloqueio: %w", err)
		}
		return true, nil
	}

	newCount, err := l.repository.IncrementCounter(l.requestKey(ip))
	if err != nil {
		return false, fmt.Errorf("erro ao incrementar contador: %w", err)
	}

	fmt.Printf("[Rate Limiter] IP %s - Uso Atualizado: %d/%d\n",
		ip, newCount, l.ipLimit)

	if newCount == 1 {
		fmt.Printf("[Rate Limiter] IP %s - Primeira requisição, definindo expiração da janela\n", ip)
		if err := l.repository.SetExpiration(l.requestKey(ip), 20); err != nil {
			return false, fmt.Errorf("erro ao definir janela: %w", err)
		}
	}

	if newCount > l.ipLimit {
		fmt.Printf("[Rate Limiter] IP %s atingiu o limite. Bloqueando por %d segundos\n",
			ip, l.blockDuration)

		if err := l.repository.SetExpiration(l.blockedKey(ip), l.blockDuration); err != nil {
			return false, fmt.Errorf("erro ao definir bloqueio: %w", err)
		}
		return true, nil
	}

	return false, nil
}

func (l *RateLimiter) IsRateLimitByToken(token string) (bool, error) {
	blocked, err := l.repository.KeyExists(l.blockedKey(token))
	fmt.Printf("[Rate Limiter] API_KEY %s - Status de Bloqueio: %v\n", token, blocked)

	if err != nil {
		return true, fmt.Errorf("erro ao verificar status de bloqueio: %w", err)
	}

	if blocked {
		fmt.Printf("[Rate Limiter] API_KEY %s está bloqueado pelos próximos %d segundos\n",
			token, l.blockDuration)
		return true, nil
	}

	currentCount, err := l.repository.GetCounter(l.requestKey(token))
	if err != nil {
		return false, fmt.Errorf("erro ao obter contagem atual: %w", err)
	}

	fmt.Printf("[Rate Limiter] API_KEY %s - Uso Atual: %d/%d\n",
		token, currentCount, l.apiKeyLimit)

	if currentCount > l.apiKeyLimit {
		fmt.Printf("[Rate Limiter] API_KEY %s excedeu o limite. Bloqueando por %d segundos\n",
			token, l.blockDuration)

		err = l.repository.SetExpiration(l.blockedKey(token), l.blockDuration)
		if err != nil {
			return false, fmt.Errorf("erro ao definir bloqueio: %w", err)
		}
		return true, nil
	}

	newCount, err := l.repository.IncrementCounter(l.requestKey(token))
	if err != nil {
		return false, fmt.Errorf("erro ao incrementar contador: %w", err)
	}

	fmt.Printf("[Rate Limiter] API_KEY %s - Uso Atualizado: %d/%d\n",
		token, newCount, l.apiKeyLimit)

	if newCount == 1 {
		fmt.Printf("[Rate Limiter] API_KEY %s - Primeira requisição, definindo expiração da janela\n", token)
		if err := l.repository.SetExpiration(l.requestKey(token), 20); err != nil {
			return false, fmt.Errorf("erro ao definir janela: %w", err)
		}
	}

	if newCount > l.apiKeyLimit {
		fmt.Printf("[Rate Limiter] API_KEY %s atingiu o limite. Bloqueando por %d segundos\n",
			token, l.blockDuration)

		if err := l.repository.SetExpiration(l.blockedKey(token), l.blockDuration); err != nil {
			return false, fmt.Errorf("erro ao definir bloqueio: %w", err)
		}
		return true, nil
	}

	return false, nil
}

func (l *RateLimiter) requestKey(identifier string) string {
	requestKey := fmt.Sprintf("req_count:%s", identifier)
	fmt.Printf("Chave de request: %s\n", requestKey)
	return requestKey
}

func (l *RateLimiter) blockedKey(identifier string) string {
	blockKey := fmt.Sprintf(fmt.Sprintf("block:%s", identifier))
	fmt.Printf("Chave de Bloqueio: %s\n", blockKey)
	return blockKey
}
