package services

import (
	"fmt"
	"time"

	"github.com/ericp/chronos-bot-reminder/internal/database"
	"github.com/ericp/chronos-bot-reminder/internal/database/models"
)

const (
	// CacheKeyAccountByDiscordID is the cache key format for account lookup by Discord ID
	CacheKeyAccountByDiscordIDFormat = "account:discord:%s"
	// CacheKeyAccountByID is the cache key format for account lookup by account ID
	CacheKeyAccountByIDFormat = "account:id:%s"
	// CacheKeyAccountIdentity is the cache key format for identity lookup
	CacheKeyAccountIdentityFormat = "identity:%s:%s"
	
	// AccountCacheDuration is how long to keep account data cached
	AccountCacheDuration = 12 * time.Hour
)

// GetAccountCacheKeyByDiscordID generates a cache key for a Discord user ID
func GetAccountCacheKeyByDiscordID(discordID string) string {
	return fmt.Sprintf(CacheKeyAccountByDiscordIDFormat, discordID)
}

// GetAccountCacheKeyByID generates a cache key for an account ID
func GetAccountCacheKeyByID(accountID string) string {
	return fmt.Sprintf(CacheKeyAccountByIDFormat, accountID)
}

// GetIdentityCacheKey generates a cache key for an identity lookup
func GetIdentityCacheKey(provider string, externalID string) string {
	return fmt.Sprintf(CacheKeyAccountIdentityFormat, provider, externalID)
}

// InvalidateAccountCache invalidates all cache entries related to an account
func InvalidateAccountCache(account *models.Account) error {
	if account == nil {
		return nil
	}

	// Delete the account ID cache
	if err := database.DeleteCache(GetAccountCacheKeyByID(account.ID.String())); err != nil {
		// Log but don't fail - cache deletion is non-critical
		fmt.Printf("[CACHE] Warning: Failed to delete account cache for ID %s: %v\n", account.ID, err)
	}

	// Delete identity caches for all associated identities
	if account.Identities != nil {
		for _, identity := range account.Identities {
			if err := database.DeleteCache(GetIdentityCacheKey(identity.Provider.String(), identity.ExternalID)); err != nil {
				// Log but don't fail
				fmt.Printf("[CACHE] Warning: Failed to delete identity cache for %s:%s: %v\n", identity.Provider, identity.ExternalID, err)
			}
		}
	}

	// Also try to delete Discord cache if we know the external ID
	if account.Identities != nil {
		for _, identity := range account.Identities {
			if identity.Provider == models.ProviderDiscord {
				if err := database.DeleteCache(GetAccountCacheKeyByDiscordID(identity.ExternalID)); err != nil {
					fmt.Printf("[CACHE] Warning: Failed to delete Discord account cache for %s: %v\n", identity.ExternalID, err)
				}
			}
		}
	}

	return nil
}

// CacheAccount stores an account in cache with its relationships
func CacheAccount(account *models.Account) error {
	if account == nil {
		return nil
	}

	// Cache by account ID
	if err := database.SetCache(GetAccountCacheKeyByID(account.ID.String()), account, AccountCacheDuration); err != nil {
		fmt.Printf("[CACHE] Warning: Failed to cache account by ID %s: %v\n", account.ID, err)
	}

	// Cache by Discord ID if applicable
	if account.Identities != nil {
		for _, identity := range account.Identities {
			if identity.Provider == models.ProviderDiscord {
				if err := database.SetCache(GetAccountCacheKeyByDiscordID(identity.ExternalID), account, AccountCacheDuration); err != nil {
					fmt.Printf("[CACHE] Warning: Failed to cache account by Discord ID %s: %v\n", identity.ExternalID, err)
				}
			}
			// Cache identity for quick lookups
			if err := database.SetCache(GetIdentityCacheKey(identity.Provider.String(), identity.ExternalID), identity, AccountCacheDuration); err != nil {
				fmt.Printf("[CACHE] Warning: Failed to cache identity %s:%s: %v\n", identity.Provider, identity.ExternalID, err)
			}
		}
	}

	return nil
}

// GetCachedAccount retrieves an account from cache by account ID
func GetCachedAccount(accountID string) (*models.Account, error) {
	var account models.Account
	err := database.GetCache(GetAccountCacheKeyByID(accountID), &account)
	if err != nil {
		// Cache miss is not an error
		return nil, nil
	}
	return &account, nil
}

// GetCachedAccountByDiscordID retrieves an account from cache by Discord ID
func GetCachedAccountByDiscordID(discordID string) (*models.Account, error) {
	var account models.Account
	err := database.GetCache(GetAccountCacheKeyByDiscordID(discordID), &account)
	if err != nil {
		// Cache miss is not an error
		return nil, nil
	}
	return &account, nil
}
