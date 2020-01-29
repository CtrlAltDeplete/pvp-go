FROM iron/go

ARG db_endpoint
ARG db_password
ARG db_user
ENV PVP_GO_DB_ENDPOINT=$db_endpoint
ENV PVP_GO_DB_PASSWORD=$db_password
ENV PVP_GO_DB_USER=$db_user

RUN mkdir /app
ADD main /app/
WORKDIR /app
CMD ["./main"]