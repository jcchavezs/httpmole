FROM scratch

ARG TARGETOS TARGETARCH

COPY build/httpmole /

EXPOSE 8081

ENTRYPOINT ["/httpmole"]
