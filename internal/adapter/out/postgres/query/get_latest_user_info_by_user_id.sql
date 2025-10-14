SELECT id, weight, height, date, age, user_id
FROM user_info 
WHERE user_id = $1
ORDER BY date DESC, id DESC
LIMIT 1