### Command
```
ansible-playbook -i hosts.ini deploy.yaml --tags "deployment"
ansible-playbook -i hosts.ini deploy.yaml --skip-tags "scheduler"
```