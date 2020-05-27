FROM buildkite/puppeteer:latest

RUN apt-get update
RUN apt-get -y install jq
RUN apt-get -y install curl
RUN apt-get -y install unzip
RUN apt-get -y install python
RUN apt-get -y install git
RUN curl "https://s3.amazonaws.com/aws-cli/awscli-bundle.zip" -o "awscli-bundle.zip" && \ 
    unzip awscli-bundle.zip && \
    ./awscli-bundle/install -i /usr/local/aws -b /usr/local/bin/aws

WORKDIR /work
ADD ./pipeline_script.sh .
ADD ./integration.js /util/

CMD ["./pipeline_script.sh"]