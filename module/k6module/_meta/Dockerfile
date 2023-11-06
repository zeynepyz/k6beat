FROM grafana/k6:latest


COPY ./script.js /home
CMD ["run", "/home/script.js"]
ENV K6_BROWSER_ENABLED=true
EXPOSE 6565

HEALTHCHECK --interval=1s --retries=90 CMD wget --spider --server-response http://localhost:6565 || exit 1