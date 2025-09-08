# Use Ubuntu 22.04 as base
FROM ubuntu:22.04

ENV DEBIAN_FRONTEND=noninteractive
ENV COMPOSER_ALLOW_SUPERUSER=1

WORKDIR /var/www/api

# Install system dependencies
RUN apt-get update -y && \
    apt-get install -y \
    git unzip libzip-dev libicu-dev libonig-dev libpng-dev \
    curl software-properties-common apt-transport-https ca-certificates gnupg g++ make nginx

# Add PHP repository
RUN add-apt-repository ppa:ondrej/php -y && apt-get update -y

# Install PHP + extensions
RUN apt-get install -y \
    php7.4 php7.4-fpm php7.4-cli php7.4-common \
    php7.4-mysql php7.4-xml php7.4-mbstring php7.4-curl php7.4-zip \
    php7.4-gd php7.4-intl php7.4-bcmath php7.4-soap php7.4-opcache

# Install Composer
RUN curl -sS https://getcomposer.org/installer | php -- --install-dir=/usr/local/bin --filename=composer

# Copy composer files first (for caching)
COPY composer.json composer.lock ./

# Install dependencies (no dev in production)
RUN composer install --no-interaction --no-ansi --no-progress --optimize-autoloader --no-dev

# Copy application
COPY . .

# Prepare Symfony cache/logs
RUN mkdir -p var/cache var/log && \
    chown -R www-data:www-data var && \
    chmod -R 775 var

# Configure Nginx to listen on Render's $PORT
RUN echo ' \
server { \
    listen ${PORT}; \
    root /var/www/api/public; \
    index index.php index.html; \
    server_name _; \
    \
    access_log /var/log/nginx/access.log; \
    error_log /var/log/nginx/error.log; \
    \
    location / { \
        try_files $uri /index.php$is_args$args; \
    } \
    \
    location ~ \.php$ { \
        include fastcgi_params; \
        fastcgi_pass unix:/var/run/php/php7.4-fpm.sock; \
        fastcgi_param SCRIPT_FILENAME $document_root$fastcgi_script_name; \
    } \
}' > /etc/nginx/sites-available/default

# Log config
RUN sed -i 's|;error_log = php_errors.log|error_log = /var/log/php7.4-fpm.log|' /etc/php/7.4/fpm/php-fpm.conf

# Symfony optimize
RUN php bin/console cache:clear --no-warmup && \
    php bin/console cache:warmup && \
    php bin/console assets:install --symlink

# Expose Render dynamic port
EXPOSE ${PORT}

# Start PHP-FPM + Nginx
CMD service php7.4-fpm start && nginx -g "daemon off;"
