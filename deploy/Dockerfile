FROM centos:latest

RUN mkdir -p /jipai/conf

VOLUME /jipai/conf

EXPOSE 4000

WORKDIR /jipai/conf

ADD assets /jipai/assets
ADD templates /jipai/templates
ADD config.json /jipai/conf/config.json

CMD /jipai/jipai -conf /jipai/conf/config.json

ADD jipai /jipai/jipai
