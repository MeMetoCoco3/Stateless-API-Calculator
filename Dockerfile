FROM debian:stable-slim
 
COPY main /bin/Outlawed
COPY middleware/middleware.go /bin/middleware



CMD ["/bin/Outlawed"]


