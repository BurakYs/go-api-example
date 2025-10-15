local key = KEYS[1]
local max = tonumber(ARGV[1])
local window = tonumber(ARGV[2])
local now = tonumber(ARGV[3])

redis.call("ZREMRANGEBYSCORE", key, 0, now - window * 1000)

local total = redis.call("ZCARD", key)
if total < max then
  redis.call("ZADD", key, now, tostring(now))
  redis.call("PEXPIRE", key, window * 1000)

  local earliest = redis.call("ZRANGE", key, 0, 0, "WITHSCORES")
  local earliestScore = tonumber(earliest[2])
  local waitMs = (earliestScore + window * 1000) - now
  if waitMs < 0 then waitMs = 0 end
  local resetSeconds = math.floor((waitMs + 999) / 1000)

  local remaining = max - total - 1
  return {remaining, resetSeconds}
else
  local earliest = redis.call("ZRANGE", key, 0, 0, "WITHSCORES")
  local earliestScore = tonumber(earliest[2])
  local waitMs = (earliestScore + window * 1000) - now
  if waitMs < 0 then waitMs = 0 end
  local resetSeconds = math.floor((waitMs + 999) / 1000)
  return {-1, resetSeconds}
end