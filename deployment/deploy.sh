pem_location=$1
ip=$2
ssh="ssh -i $pem_location ec2-user@${ip}"
release_dir=deployment/releases/$(date "+%F-%T")

rm -rf $release_dir
make docker_server

mkdir $release_dir

docker save -o ${release_dir}/qode.tar qode

scp -r -i $pem_location $release_dir/qode.tar ec2-user@${ip}:/home/ec2-user/release/
scp -r -i $pem_location deployment/qode/migrations/* ec2-user@${ip}:/home/ec2-user/release/migrations/
${ssh} "docker load -i ./release/qode.tar"
${ssh} "pushd release/deployment && docker-compose down && docker-compose up -d"

