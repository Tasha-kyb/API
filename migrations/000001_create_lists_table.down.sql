-- Откат миграции: удаление таблицы
DROP INDEX IF EXISTS idx_lists_created_at;
DROP TABLE IF EXISTS lists;