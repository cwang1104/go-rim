local topicZSet = KEYS[1]
local topicHash = KEYS[2]
local key = ARGV[1]


redis.call("zrem",topicZSet,key)
redis.call("hdel",topicHash,key)

return 1