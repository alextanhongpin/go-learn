Minimal golang binary with upx:

```
ARG upx_version=5.1.0
ARG TARGETARCH=amd64

RUN update-ca-certificates && apt update && apt install -y tzdata fontconfig && \
	curl -Ls https://github.com/upx/upx/releases/download/v${upx_version}/upx-${upx_version}-${TARGETARCH}_linux.tar.xz -o - | tar xvJf - -C /tmp && \
  cp /tmp/upx-${upx_version}-${TARGETARCH}_linux/upx /usr/local/bin/ && \
  chmod +x /usr/local/bin/upx && \
	apt autoremove && apt clean && \
 	rm -rf /var/lib/apt/lists/* /tmp/upx-*
WORKDIR /root/
RUN find ./bin -name '*.so' -exec chmod +x {} \; && find ./bin -name '*.so' -exec upx -9 -k {} \;
```

```
curl -Ls https://github.com/upx/upx/releases/download/v5.1.0/upx-5.1.0-amd64_linux.tar.xz -o - | tar xvJf - -C /tmp && \
  cp /tmp/upx-5.1.0-amd64_linux/upx /usr/local/bin/ && \
  chmod +x /usr/local/bin/upx && \
	apt autoremove && apt clean && \
 	rm -rf /var/lib/apt/lists/* /tmp/upx-*
```
