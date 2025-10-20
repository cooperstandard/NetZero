# TODO

* on startup, create database if not exist and use goose to run up migrations
* pay debt / settle up
* leave group
* patch user (and the other entities which may need patching)
* delete transaction
* functional testing (with option to run through a local sqlite instance or through psql db)
* bruno for balances and transactions
* add limit to group size, important because well connected groups of size n have n^2 balance records which could cause db lockups
* settleup
* document how to configure and run app with readme and sample .env, use docker compose secret file
* github action to automatically create incremented tags when PR merged to main with "chore:", "fix:", "patch:", or "feat:" in the title
