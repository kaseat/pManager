#!/bin/bash
cat cleanup.sql \
functions/pseudo_encrypt_24.sql \
tables/user_roles.sql \
tables/users.sql \
tables/sync_providers.sql \
tables/user_sync.sql \
tables/portfolios.sql \
tables/currencies.sql \
tables/securities_types.sql \
tables/exchange.sql \
tables/securities.sql \
tables/operation_types.sql \
tables/operations.sql \
tables/prices.sql \
tables/settings.sql \
post_deployment.sql > res.sql