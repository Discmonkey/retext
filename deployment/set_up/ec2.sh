# run this file from the project root
pem_location=$1
ip=$2
release_name=$(date "+%F-%T")
release_dir=deployment/releases/${release_name}

ssh="ssh -i ${pem_location} ec2-user@${ip}"

${ssh} "sudo yum update -y"
${ssh} "sudo amazon-linux-extras install docker"
${ssh} "sudo service docker start"
${ssh} "sudo usermod -a -G docker ec2-user"

echo "checking successful provision"
${ssh} "docker info"

${ssh} "sudo curl -L https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m) -o /usr/local/bin/docker-compose"
${ssh} "sudo chmod +x /usr/local/bin/docker-compose"

rm -rf $release_dir ${release_dir}.tar
make all
mkdir $release_dir

docker save -o ${release_dir}/qode.tar qode
docker save -o ${release_dir}/qode_db_loader.tar qode_db_loader
docker save -o ${release_dir}/postgres.tar postgres

cp -r deployment/qode ${release_dir}/deployment
mv ${release_dir}/deployment/docker-compose.yml.release ${release_dir}/deployment/docker-compose.yml

scp -r -i $pem_location $release_dir ec2-user@${ip}:/home/ec2-user/release

${ssh} "docker load -i release/${release_name}/qode.tar"
${ssh} "docker load -i release/${release_name}/qode_db_loader.tar"
${ssh} "docker load -i release/${release_name}/postgres.tar"

${ssh} "cd release/${release_name}/deployment && docker-compose up -d"