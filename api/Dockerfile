# Use official PHP 7.4 FPM image
FROM php:7.4-fpm

# Set working directories
WORKDIR /app

# Install system dependencies
RUN apt-get update -y && apt-get install -y \
    git \
    unzip \
    libzip-dev \
    libicu-dev \
    libonig-dev \
    libpng-dev \
    curl \
    nginx \
    supervisor \
    make \
    gcc \
    libc6-dev

# Install PHP extensions
RUN docker-php-ext-install \
    pdo_mysql \
    intl \
    mbstring \
    zip \
    gd \
    bcmath \
    soap \
    opcache

# Install Composer properly (Fix for Railway deployment)
RUN curl -sS https://getcomposer.org/installer | php -- --install-dir=/usr/local/bin --filename=composer

# Verify Composer installation
RUN composer --version

# Create necessary directories
RUN mkdir -p /var/log/nginx /var/cache/nginx

# Copy application files
COPY . /app/

# Ensure Composer runs from correct path
ENV PATH="/usr/local/bin:$PATH"

# Install dependencies (Fix for Composer issue)
RUN cd /app && composer install --no-dev --optimize-autoloader --no-interaction --ignore-platform-reqs

# Compile Swiss Ephemeris if needed
RUN if [ -d "/app/swetest/src" ]; then \
    cd /app/swetest/src && make clean && make && chmod 777 . && chmod +x swetest; \
    fi

# Set permissions
RUN chown -R www-data:www-data /app

# Configure Nginx
RUN echo ' \
server { \
    listen 9393; \
    root /app/public; \
    index index.php index.html index.htm; \
    server_name _; \
    \
    access_log /var/log/nginx/access.log; \
    error_log /var/log/nginx/error.log; \
    \
    location / { \
        try_files $uri $uri/ /index.php?$query_string; \
    } \
    \
    location ~ \.php$ { \
        fastcgi_pass unix:/run/php/php7.4-fpm.sock; \
        fastcgi_index index.php; \
        include fastcgi_params; \
        fastcgi_param SCRIPT_FILENAME $document_root$fastcgi_script_name; \
    } \
}' > /etc/nginx/sites-available/default

# Configure Supervisor
RUN echo '[supervisord] \
nodaemon=true \
[program:nginx] \
command=/usr/sbin/nginx -g "daemon off;" \
autostart=true \
autorestart=true \
stderr_logfile=/var/log/nginx/error.log \
stdout_logfile=/var/log/nginx/access.log \
[program:php-fpm] \
command=docker-php-entrypoint php-fpm \
autostart=true \
autorestart=true \
stderr_logfile=/var/log/php7.4-fpm.log \
stdout_logfile=/var/log/php7.4-fpm.log' > /etc/supervisor/conf.d/supervisord.conf

# Expose port
EXPOSE 9393

# Start Supervisor
CMD ["supervisord", "-c", "/etc/supervisor/conf.d/supervisord.conf"]