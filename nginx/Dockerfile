FROM nginx:1.26.1-alpine3.19

COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 8291

CMD ["nginx", "-g", "daemon off;"]