FROM postgres:alpine

WORKDIR /app

# Create the directory for PostgreSQL data
RUN mkdir -p /var/lib/postgresql/data

# Initialize PostgreSQL data directory and set permissions
RUN chown -R postgres:postgres /var/lib/postgresql/data && \
    chmod 700 /var/lib/postgresql/data

# Set environment variables
ENV POSTGRES_DB=autotransport
ENV POSTGRES_USER=postgres
ENV POSTGRES_PASSWORD=example

# Switch to postgres user to initialize database
USER postgres

# Initialize PostgreSQL data directory
RUN initdb -D /var/lib/postgresql/data

# Copy your application files
COPY ./bin/autotransport /app/autotransport
COPY ./static /app/static
COPY ./ddl.sql /app/ddl.sql

# Expose ports
EXPOSE 8085
EXPOSE 5432

# Start PostgreSQL and then run the application
# CMD ["/app/docker-entrypoint.sh"]
