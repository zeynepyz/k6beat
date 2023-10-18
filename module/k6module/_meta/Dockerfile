ARG K6_VERSION=0.46.0
FROM loadimpact/k6:${K6_VERSION}

COPY ./script.js /home

CMD ["run", "/home/script.js"]

HEALTHCHECK --interval=1s --retries=90 CMD wget --spider --server-response http://localhost:6565 || exit 1