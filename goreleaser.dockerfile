FROM scratch
COPY op-plasma-sunrise /usr/bin/op-plasma-sunrise
ENTRYPOINT ["/usr/bin/op-plasma-sunrise"]