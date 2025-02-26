# Use official PHP 7.4 FPM image (lighter and optimized)
FROM php:7.4-fpm

# Set working directory
WORKDIR /var/www/api

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
    supervisor

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

# Install Composer
RUN curl -sS https://getcomposer.org/installer | php -- --install-dir=/usr/local/bin --filename=composer && \
    /usr/local/bin/composer --version

# Copy application files
COPY . .

# Install dependencies with full Composer path
RUN /usr/local/bin/composer install --no-dev --optimize-autoloader --no-interaction

# Set proper permissions
RUN mkdir -p /var/www/api/var/cache /var/www/api/var/log && \
    chown -R www-data:www-data /var/www/api/var && \
    chmod -R 777 /var/www/api/var

# Configure Nginx
RUN echo ' \
server { \
    listen 9393; \
    root /var/www/api/public; \
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

# Configure Supervisor to manage processes
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

# Start Supervisor to run both Nginx and PHP-FPM
CMD ["supervisord", "-c", "/etc/supervisor/conf.d/supervisord.conf"]
