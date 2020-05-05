FROM gcr.io/distroless/static
LABEL Brayden Winterton <brayden_winterton@byu.edu>

COPY av-uapi av-uapi

ENTRYPOINT ["/av-uapi"]