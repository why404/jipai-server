IMAGE="docker.cn/googollee/pili-jipai:v1"
BINARY="jipai"

cp -r templates assets deploy && docker build $1 -t output . && docker run --rm -v `pwd`:/go/src/$BINARY output go build -o /go/src/$BINARY/deploy/$BINARY $BINARY && docker build -t "$IMAGE" deploy/ && docker push "$IMAGE" && rm -r deploy/$BINARY
