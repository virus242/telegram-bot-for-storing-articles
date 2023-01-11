FROM golang

WORKDIR /github.com/virus242/telegram-bot-for-storing-articles/

COPY . .  

RUN go mod download

EXPOSE 3000

CMD [ "go", "run", "cmd/main.go", "--TOKEN=5502453090:AAHDWYzON-7jEC0lZ698l3FM0q29AB4soPs" ]



