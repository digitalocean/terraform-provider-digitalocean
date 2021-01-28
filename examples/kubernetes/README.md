terraform apply -var="write_kubeconfig=true"
export KUBECONFIG=$(terraform output -raw kubeconfig_path