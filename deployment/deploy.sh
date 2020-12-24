pem_location=$1
ip=$2
release_dir=deployment/releases/$(date -I)
ssh="ssh -i $pem_location ec2-user@18.223.164.85"

make docker_server
mkdir $release_dir

docker save -o ${release_dir}/qode.tar qode

scp -i $pem_location $release_dir/qode.tar ec2-user@${ip}:/home/ec2-user/

${ssh} "docker load -i qode.tar"
${ssh} "pushd retext/deployment/qode && docker-compose down && docker-compose up -d"

