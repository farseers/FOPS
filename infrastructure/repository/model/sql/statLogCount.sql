SELECT app_name,
       log_level,
       COUNT(*) AS log_count
FROM linkTrace.log_data
WHERE create_at >= (NOW() - INTERVAL 30 MINUTE)
GROUP BY app_name, log_level;