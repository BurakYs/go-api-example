local key = KEYS[1]
local max = tonumber(ARGV[1])
local window = tonumber(ARGV[2])
local now = tonumber(ARGV[3])

local current = tonumber(redis.call("GET", key) or "0")
if current >= max then
  local ttl = redis.call("TTL", key)
  if ttl < 0 then
    redis.call("EXPIRE", key, window)
    ttl = window
  end

  if ttl < 0 then ttl = 0 end
  return {-1, ttl}
end

current = redis.call("INCR", key)
if current == 1 then
  redis.call("EXPIRE", key, window)
end

local ttl = redis.call("TTL", key)
if ttl < 0 then ttl = 0 end
local remaining = max - current
return {remaining, ttl}