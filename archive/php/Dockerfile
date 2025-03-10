FROM composer:latest AS builder

COPY composer.json composer.lock ./

WORKDIR /app

RUN composer install

COPY . . 



FROM php:8.3.12-apache

WORKDIR /var/www/html/

RUN docker-php-ext-install mysqli

RUN a2enmod rewrite

COPY --from=builder /app /var/www/html/

EXPOSE 80
EXPOSE 443

CMD ["apache2-foreground"]