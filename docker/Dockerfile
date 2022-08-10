FROM debian:latest

# install and config zsh
# create non-root user
RUN apt-get update &&\
    apt-get install -y zsh zsh-syntax-highlighting zsh-autosuggestions &&\
    chsh -s /bin/zsh root &&\
    apt-get install -y sudo &&\ 
    groupadd -g 1000 user &&\
    useradd -rm -d /home/user -s /bin/zsh -u 1000 -g 1000 -G sudo user &&\
    echo 'user:user' | chpasswd
ADD .zshrc /root/.zshrc
ADD .zshrc /home/user/.zshrc

# download kafka
RUN apt-get install -y wget openjdk-11-jdk &&\
    cd /home/user &&\
    wget https://dlcdn.apache.org/kafka/3.2.1/kafka_2.13-3.2.1.tgz &&\
    tar -xzf kafka_2.13-3.2.1.tgz && rm kafka_2.13-3.2.1.tgz &&\
    chown user:user -R kafka_2.13-3.2.1 &&\
    echo "\n\nPATH=\$PATH:$(pwd)/kafka_2.13-3.2.1/bin" >> .zshrc

# install common package
RUN apt-get install -y neovim

# change default user
WORKDIR /home/user/
USER user

# change default shell
CMD [ "/bin/zsh" ]