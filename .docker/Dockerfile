FROM debian:latest

# install and config zsh
RUN apt-get update &&\
    apt-get install -y zsh zsh-syntax-highlighting zsh-autosuggestions &&\
    chsh -s /bin/zsh root
ADD .zsh_history /root/.zsh_history
ADD .zshrc /root/.zshrc

# download kafka
RUN apt-get install -y wget openjdk-11-jdk &&\
    cd /root &&\
    wget https://dlcdn.apache.org/kafka/3.2.1/kafka_2.13-3.2.1.tgz &&\
    tar -xzf kafka_2.13-3.2.1.tgz &&\
    rm kafka_2.13-3.2.1.tgz &&\
    echo "\n\nPATH=\$PATH:$(pwd)/kafka_2.13-3.2.1/bin" >> .zshrc

# install common package
RUN apt-get install -y neovim

# change default work directory
WORKDIR /root

# expose port
EXPOSE 9090-9099

# run kafka-zookeeper by default
CMD [ \
    "/root/kafka_2.13-3.2.1/bin/zookeeper-server-start.sh", \
    "/root/kafka_2.13-3.2.1/config/zookeeper.properties" \
    ]