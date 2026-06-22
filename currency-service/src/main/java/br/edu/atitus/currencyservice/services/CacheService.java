package br.edu.atitus.currencyservice.services;

import org.springframework.cache.Cache;
import org.springframework.cache.CacheManager;
import org.springframework.stereotype.Service;

import java.math.BigDecimal;

@Service
public class CacheService {
    private final CacheManager cacheManager;

    public CacheService(CacheManager cacheManager) {
        this.cacheManager = cacheManager;
    }

    public BigDecimal get(String cacheName, String key) {

        Cache cache = cacheManager.getCache(cacheName);

        if (cache == null) {
            return null;
        }

        Cache.ValueWrapper wrapper = cache.get(key);

        if (wrapper == null) {
            return null;
        }

        return (BigDecimal) wrapper.get();
    }

    public void put(String cacheName, String key, Object value) {

        Cache cache = cacheManager.getCache(cacheName);

        if (cache != null) {
            cache.put(key, value);
        }
    }
}
