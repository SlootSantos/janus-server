FROM buildkite/puppeteer:latest

# install required tools
RUN apt-get update
RUN apt-get -y install jq curl unzip python git

# add non-root user for execution
RUN useradd -u 8899 pipeline
RUN mkdir -p /home/pipeline/app /home/pipeline/source
RUN touch /home/pipeline/.bashrc
RUN chown -R pipeline /home/pipeline

# install aws CLI
# TODO: should change /usr/local/bin/aws installation path
RUN curl "https://s3.amazonaws.com/aws-cli/awscli-bundle.zip" -o "awscli-bundle.zip" && \ 
    unzip awscli-bundle.zip && \
    ./awscli-bundle/install -i /usr/local/aws -b /usr/local/bin/aws 

# mv preinstalled npm configuration to pipeline's home
RUN mv ~/.npm /home/pipeline/.npm
RUN mv ~/.config /home/pipeline/.config
RUN chown -R pipeline /home/pipeline/.npm /home/pipeline/.config

WORKDIR /home/pipeline/app

# copy pipeline assets
COPY ./pipeline_script.sh .
COPY ./integration.js /home/pipeline/util/
RUN chmod 555 pipeline_script.sh
RUN chmod 555 /home/pipeline/util/integration.js

# do everything as pipeline from now on
USER pipeline

# install NVM as pipeline 
RUN curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.35.3/install.sh | bash

CMD ["./pipeline_script.sh"]