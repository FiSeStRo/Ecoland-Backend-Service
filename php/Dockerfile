FROM php:8.3.12-apache

COPY . /var/www/html/

RUN docker-php-ext-install mysqli

RUN a2enmod rewrite

EXPOSE 80
EXPOSE 443

CMD ["apache2-foreground"]
