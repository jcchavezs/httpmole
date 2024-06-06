FROM --platform=${TARGETPLATFORM} scratch

ARG TARGETOS TARGETARCH

COPY build/httpmole-${TARGETOS}-${TARGETARCH} /httpmole

EXPOSE 10080

ENTRYPOINT ["/httpmole", "-p", "10080"]
