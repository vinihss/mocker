FROM nginx:alpine

# Copy configs and files
COPY nginx.conf /etc/nginx/nginx.conf
COPY frontend/dist /usr/share/nginx/html/
COPY mocker /mocker

# Create startup script
RUN echo '#!/bin/sh' > /start.sh && \
    echo '/mocker &' >> /start.sh && \
    echo 'nginx -g "daemon off;"' >> /start.sh && \
    chmod +x /start.sh

EXPOSE 80 8080

CMD ["/start.sh"]