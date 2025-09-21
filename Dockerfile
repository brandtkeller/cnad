FROM gcr.io/distroless/base-debian12:latest
COPY myapp /myapp
EXPOSE 8080
ENTRYPOINT ["/myapp"]
