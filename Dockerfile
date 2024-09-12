FROM --platform=${TARGETPLATFORM} scratch

ARG TARGETOS TARGETARCH

COPY build/httpmole-${TARGETOS}-${TARGETARCH} /httpmole

ENTRYPOINT ["/httpmole"]