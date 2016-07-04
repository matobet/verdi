-- update-indexes TYPE ID field value [field value...]
local type = KEYS[1]
local id = KEYS[2]
local newIndexValues = {}
local indexes = {}
local nextKey
for i, v in ipairs(ARGV) do
    if i % 2 == 1 then
        table.insert(indexes, v)
        nextKey = v
    else
        newIndexValues[nextKey] = v
    end
end
local old = redis.call('HMGET', string.format('%s:%s', type, id), unpack(indexes))
local updates = 0
for i, index in ipairs(indexes) do
    if old[i] and old[i] ~= newIndexValues[index] then
        redis.call('DEL', string.format('q:%s:%s:%s', type, index, old[i]))
    end
    if (not old[i] or old[i] ~= newIndexValues[index]) and newIndexValues[index] ~= '' then
        redis.call('SET', string.format('q:%s:%s:%s', type, index, newIndexValues[index]), id)
        updates = updates + 1
    end
end
return updates
