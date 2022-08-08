#! /bin/zsh

readonly ARCH=$(uname -m)
echo "Found Arch to be $ARCH"

curl "https://awscli.amazonaws.com/awscli-exe-linux-$ARCH-2.4.17.zip" -o "awscliv2.zip"
unzip awscliv2.zip
sudo ./aws/install
rm -rf awscliv2.zip ./aws
