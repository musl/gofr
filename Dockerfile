FROM scratch
ADD cmd/gofrd/gofrd /gofrd
ADD cmd/gofrd/build /build
EXPOSE 8080
CMD ["/gofrd"]
