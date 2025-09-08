# ======================
# 1. Base image
# ======================
FROM php:7.4-fpm

# Install system dependencies & PHP extensions
RUN apt-get update && apt-get install -y \
    git unzip curl libicu-dev libpq-dev libzip-dev libxml2-dev libonig-dev \
    && docker-php-ext-install intl pdo pdo_mysql zip opcache bcmath soap \
    && docker-php-ext-enable opcache

# ======================
# 2. Composer
# ======================
COPY --from=composer:2 /usr/bin/composer /usr/bin/composer

# ======================
# 3. App code
# ======================
WORKDIR /var/www/html

# Copy composer manifests first (for caching)
COPY composer.json composer.lock* ./

# Install PHP dependencies (optimized autoloader, no dev)
RUN composer install --no-interaction --no-ansi --no-progress --optimize-autoloader --no-dev

# Copy rest of the application
COPY . .

# ======================
# 4. Permissions
# ======================
RUN chown -R www-data:www-data /var/www/html \
    && chmod -R 755 /var/www/html

# ======================
# 5. PHP-FPM config
# ======================
EXPOSE 9000
CMD ["php-fpm"]
