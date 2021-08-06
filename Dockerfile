FROM public.ecr.aws/lambda/provided:al2 as build
RUN yum install -y golang
RUN go env -w GOPROXY=direct
RUN go get -u github.com/fogleman/primitive
ADD go.mod go.sum ./
RUN go mod download
ADD . .
RUN go build -o /main
FROM public.ecr.aws/lambda/provided:al2
COPY --from=build /main /main
COPY --from=build /root/go/bin/primitive /primitive
COPY entry.sh /
RUN chmod 755 /entry.sh
ENTRYPOINT [ "/entry.sh" ]
