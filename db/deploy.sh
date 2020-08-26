#!/bin/bash
cat cleanup.sql \
functions/pseudo_encrypt_24.sql \
tables/user_roles.sql \
tables/users.sql \
tables/portfolios.sql \
tables/currencies.sql \
tables/securities_types.sql \
tables/securities.sql \
tables/operation_types.sql \
tables/operations.sql \
tables/prices.sql \
post_deployment.sql > res.sql