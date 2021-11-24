FROM mariadb:10.5

# ENV MARIADB_ALLOW_EMPTY_ROOT_PASSWORD yes
ENV MARIADB_ROOT_PASSWORD rootpassword
# COPY setup.sh /mysql/setup.sh
COPY ./sql/init.sql /docker-entrypoint-initdb.d
# COPY ./sql/privileges.sql /docker-entrypoint-initdb.d


# RUN sh /mysql/setup.sh`