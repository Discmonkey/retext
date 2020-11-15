release_dir=deployment/releases/$(date -I)
pem_location=/home/max/sshkeys/amazon.pem
ssh="ssh -i $pem_location ec2-user@18.223.164.85"

make docker_server
mkdir $release_dir

docker save -o ${release_dir}/qode.tar qode

scp -i $pem_location $release_dir/qode.tar ec2-user@18.223.164.85:/home/ec2-user/

${ssh} "docker load -i qode.tar"
${ssh} "pushd retext/deployment/qode && ls && docker-compose down --volumes && docker-compose up -d"

