-- List activitiesc
select user_id, activity, starttime 
FROM individual_activities 
INNER JOIN users 
    ON user_id = users.id 
    WHERE line_id = 'U2809deeb91fd72aac99184b430865922' 
ORDER BY starttime DESC;
