FROM public.ecr.aws/lambda/ruby:2

COPY app.rb ./

CMD [ "app.handler" ]
