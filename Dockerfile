FROM vault:0.8.3
ADD bin/linux/actions-plugin /vault-plugins/actions-plugin
ADD bin/linux/checksum /vault-plugins/checksum 
ENV VAULT_DEV_ROOT_TOKEN_ID=1234
ENV VAULT_ADDR=http://127.0.0.1:8200
ENV VAULT_LOCAL_CONFIG='{"plugin_directory":"/vault-plugins"}'

COPY docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh
ENTRYPOINT ["docker-entrypoint.sh"]

CMD ["server", "-dev"]