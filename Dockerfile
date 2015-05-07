FROM busybox:latest
ADD scaledowntargetemulator /bin/
CMD ["/bin/scaledowntargetemulator"]
