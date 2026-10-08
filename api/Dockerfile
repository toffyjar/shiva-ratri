# Base image: Ubuntu 24.04 ships PHP 8.3 in its main archive (no PPA needed).
# PHP 7.4 is end-of-life, and current releases of this app's dependencies
# (twig, twig/extra-bundle) no longer install on it.
FROM ubuntu:24.04

ENV DEBIAN_FRONTEND=noninteractive
ENV COMPOSER_ALLOW_SUPERUSER=1
WORKDIR /var/www/api

# System dependencies + PHP 8.3 + extensions
RUN apt-get update -y && \
    apt-get install -y --no-install-recommends \
    ca-certificates curl git unzip nginx \
    php8.3-fpm php8.3-cli php8.3-common \
    php8.3-mysql php8.3-xml php8.3-mbstring php8.3-curl php8.3-zip \
    php8.3-gd php8.3-intl php8.3-bcmath php8.3-soap php8.3-opcache && \
    rm -rf /var/lib/apt/lists/*

# Install Composer
RUN curl -sS https://getcomposer.org/installer | php -- --install-dir=/usr/local/bin --filename=composer

# Install dependencies. No "|| true": if this fails the build must fail, so Render
# keeps serving the previous deploy instead of starting an image without vendor/.
COPY composer.json symfony.lock ./
RUN composer install \
    --no-interaction \
    --no-ansi \
    --no-progress \
    --optimize-autoloader \
    --no-dev \
    --no-scripts

# Copy full application and rebuild the autoloader against it
COPY . .
RUN composer dump-autoload --optimize --no-dev

# Warm the Symfony cache, then hand var/ to the PHP-FPM user
RUN php bin/console cache:clear --no-warmup && \
    php bin/console cache:warmup && \
    mkdir -p var/cache var/log && \
    chown -R www-data:www-data var && chmod -R 775 var

# Let PHP-FPM see the container's environment variables (SR_API_KEYS, SR_ALLOWED_IPS, ...): by default
# FPM clears them before running PHP.
RUN sed -i 's/^;\?clear_env = .*/clear_env = no/' /etc/php/8.3/fpm/pool.d/www.conf && \
    grep -q '^clear_env = no' /etc/php/8.3/fpm/pool.d/www.conf

# Configure Nginx
COPY docker/nginx.conf /etc/nginx/sites-available/default

EXPOSE 8080

# Start PHP-FPM + Nginx
CMD service php8.3-fpm start && nginx -g "daemon off;"
