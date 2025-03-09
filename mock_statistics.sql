-- Создаем временную таблицу с бизнесами (фиксированный business_id)
CREATE
TEMPORARY TABLE IF NOT EXISTS businesses AS
SELECT '8594581b-ee01-4860-8d61-4f07677de408' AS business_id;

-- Создаем временную таблицу с программами лояльности (1-3 программы на бизнес)
CREATE
TEMPORARY TABLE IF NOT EXISTS loyalty_programs AS
SELECT business_id,
       arrayJoin(arrayMap(x - > generateUUIDv4(), range(1 + rand() % 3))) AS program_id
FROM businesses;

-- Создаем базовый набор пользователей (1000-2000 начальных пользователей на бизнес)
CREATE
TEMPORARY TABLE IF NOT EXISTS initial_users AS
SELECT business_id,
       arrayJoin(arrayMap(x - > generateUUIDv4(), range(1000 + rand() % 1001))) AS user_id,
       toDate('2025-03-03') - INTERVAL 90 DAY + toIntervalDay(rand() % 30) as registration_date
FROM businesses;

-- Создаем дополнительных пользователей (5-15 новых пользователей каждый день после первых 30 дней)
CREATE
TEMPORARY TABLE IF NOT EXISTS additional_users AS
SELECT business_id,
       arrayJoin(arrayMap(x - > generateUUIDv4(), range(5 + rand() % 11))) AS user_id, date as registration_date
FROM businesses
    CROSS JOIN (
    SELECT arrayJoin(arrayMap(x -> toDate('2025-03-03') - INTERVAL 60 DAY + toIntervalDay(x), range (60))) as date
    ) AS dates;

-- Объединяем всех пользователей
CREATE
TEMPORARY TABLE IF NOT EXISTS users AS
SELECT *
FROM initial_users
UNION ALL
SELECT *
FROM additional_users;

-- Генерируем данные за последние 90 дней
INSERT INTO coin_balance_changes
WITH
    -- Создаем набор пользователей с их активностью
    user_activities AS (SELECT u.*,
                               -- 0-3 сканирования в день после регистрации
                               arrayJoin(arrayMap(x - > 1, range(rand() % 4))) as activity_count
                        FROM users u
                        -- 30% вероятность активности в конкретный день
                        WHERE rand() % 100 < 30
    )
SELECT b.business_id,
       u.user_id,
       lp.program_id,
       CASE
           -- QR сканирования (всегда +1 коин)
           WHEN rand() % 100 < 85 THEN toInt64(1) -- 85% случаев, +1 коин
       -- Покупки купонов (небольшие списания)
           ELSE toInt64(-(5 + rand() % 16)) -- 15% случаев, -5 до -20 монет
           END AS balance_change,
       CASE
           WHEN balance_change > 0 THEN 'qr_scan'
           ELSE 'coupon_buy'
           END AS reason,
       CASE
           WHEN reason = 'coupon_buy' THEN generateUUIDv4()
           ELSE NULL
           END AS coupon_id,
       -- Генерируем временные метки с учетом даты регистрации пользователя
       toDateTime64(
               toDateTime(u.registration_date) +
                   INTERVAL (rand() % dateDiff('day', u.registration_date, toDate('2025-03-03'))) DAY +
                   -- Больше активности в рабочие часы
               INTERVAL (
                   CASE
                       WHEN rand() % 100 < 70 THEN 9 + rand() % 9 -- 70% в рабочие часы (9-17)
                       ELSE 17 + rand() % 7 -- 30% вечером (17-23)
                       END
                   ) HOUR +
               INTERVAL (rand() % 60) MINUTE,
               3, 'UTC'
       ) AS timestamp
FROM
    businesses b
    CROSS JOIN user_activities u
    LEFT JOIN loyalty_programs lp
ON b.business_id = lp.business_id
WHERE toDateTime(u.registration_date) +
    INTERVAL (rand() % dateDiff('day'
    , u.registration_date
    , toDate('2025-03-03'))) DAY +
    INTERVAL (
    CASE
    WHEN rand() % 100
    < 70 THEN 9 + rand() % 9
    ELSE 17 + rand() % 7
    END
    ) HOUR +
    INTERVAL (rand() % 60) MINUTE
    > toDateTime(u.registration_date) -- Проверяем, что активность после регистрации
ORDER BY timestamp;

-- Очистка временных таблиц
DROP TABLE IF EXISTS businesses;
DROP TABLE IF EXISTS loyalty_programs;
DROP TABLE IF EXISTS initial_users;
DROP TABLE IF EXISTS additional_users;
DROP TABLE IF EXISTS users;