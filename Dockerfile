FROM alpine:3.7

COPY httpmole /

EXPOSE 8081

ENTRYPOINT ["/httpmole"]
