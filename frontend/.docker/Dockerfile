FROM node:14-alpine as builder
WORKDIR /app
COPY . .
RUN yarn install --frozen-lockfile
RUN yarn build

FROM nginx:alpine
COPY --from=builder /app/build /usr/share/nginx/html
