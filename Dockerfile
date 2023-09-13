FROM ubuntu:latest

RUN apt-get update --fix-missing -y
RUN apt -y update
RUN apt -y upgrade
RUN apt install -y curl
RUN apt install -y ffmpeg
RUN apt install -y golang
RUN apt install -y python3

# Took from ytdlp wiki: https://github.com/yt-dlp/yt-dlp/wiki/Installation
RUN curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -o /usr/local/bin/yt-dlp
RUN chmod 700 /usr/local/bin/yt-dlp  # Make executable

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY *.go ./
COPY enc/ ./enc

RUN go build -o /ndmb

ENTRYPOINT ["/ndmb"]
CMD ["--ytdlp", "/usr/local/bin/yt-dlp", "--ffmpeg", "/usr/bin/ffmpeg"]
