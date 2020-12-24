# run this file from the project root
pem_location=$1
ip=$2

ssh="ssh -i ${pem_location} ec2-user@${ip}"

${ssh} "sudo yum update -y"
${ssh} "sudo amazon-linux-extras install docker"
${ssh} "sudo service docker start"
${ssh} "sudo usermod -a -G docker ec2-user"

echo "checking successful provision"
${ssh} "docker info"