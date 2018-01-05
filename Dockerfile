#base image
FROM busybox:glibc

#copy
COPY ./webhook /
COPY ./pki /pki

#setting working directory
WORKDIR /
#expose port
EXPOSE 8080

#entry point
#ENTRYPOINT ["/webhook"]
CMD ["/webhook","serve"]