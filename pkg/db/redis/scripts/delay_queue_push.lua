local topicZSet = KEYS[1]
local topicHash = KEYS[2]
local key = ARGV[1]
local body = ARGV[2]
local readyTime = tonumber(ARGV[3])

-- 添加readyTime到zset
local count = redis.call("zadd", topicZSet, readyTime, key)
-- 消息已经存在
if count == 0 then
    return 0
end
-- 添加body到hash
redis.call("hsetnx", topicHash, key, body)

return 1