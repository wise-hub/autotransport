FROM postgres:alpine

WORKDIR /app

RUN mkdir -p /var/lib/postgresql/data

RUN chown -R postgres:postgres /var/lib/postgresql/data && \
    chmod 700 /var/lib/postgresql/data

ENV POSTGRES_DB=autotransport
ENV POSTGRES_USER=postgres
ENV POSTGRES_PASSWORD=example

USER postgres

RUN initdb -D /var/lib/postgresql/data

COPY ./bin/autotransport /app/autotransport
COPY ./web /app/web
COPY ./ddl.sql /app/ddl.sql

EXPOSE 8085
EXPOSE 5432