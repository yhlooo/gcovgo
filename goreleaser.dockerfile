FROM --platform=${TARGETPLATFORM} busybox:latest
COPY gcovgo /bin/gcovgo
ENTRYPOINT ["/bin/gcovgo"]
