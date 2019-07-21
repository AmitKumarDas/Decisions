```yaml
env: &env
  FLASK_APP: app.py
  FLASK_ENV: development
dependencies:
  - postgres
  - db-migrate
tasks:
- name: db-migrate
  args: [python3, migrate_database.py]
  description: Migrates the database to the proper schema
  env: *env
```
