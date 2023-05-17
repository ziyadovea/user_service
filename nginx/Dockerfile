# Official Nginx image
FROM nginx:1.23-alpine

# Copy Nginx configuration file to the container
# nginx_lb_active_active.conf or nginx_lb_active_passive.conf
COPY nginx_lb_active_active.conf /etc/nginx/conf.d

# Expose port 80
EXPOSE 80