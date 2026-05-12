FROM docker.io/bitnami/nginx:latest

USER root

RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo "Asia/Shanghai" > /etc/timezone

WORKDIR /app

RUN mkdir -p /app/client_downloads \
    && mkdir -p /app/config \
    && mkdir -p /opt/chprobe/etc \
    && chmod 777 /app/client_downloads

COPY docker/nginx.conf /opt/bitnami/nginx/conf/server_blocks/default.conf

COPY docker/start.sh /app/start.sh
RUN chmod +x /app/start.sh

COPY chprobe_web/dist /app/web

COPY chprobe_server/bin/chprobe_server /app/

COPY chprobe_client/bin/chprobe_client /app/client_downloads/

COPY chprobe_server/config/config.yaml /app/config/

COPY chprobe_server/utils/public.pem /opt/chprobe/etc/public.pem

EXPOSE 8080 32000

ENTRYPOINT []
CMD ["/app/start.sh"]
