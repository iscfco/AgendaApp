import subprocess

def execute_docker_stop(remove_volumes=False):
    """Run all tasks required to stop docker compose services"""

    print("📋 [Stop] Stopping services")

    # Docker Compose Command
    docker_cmd = 'docker-compose -f .\\artifacts\\docker-compose.yaml down'

    if remove_volumes:
        docker_cmd += ' --volumes'
    
    #check=True will raise an error if the command fails
    subprocess.run(docker_cmd, shell=True, check=True)

def  execute_docker_force_stop():
    execute_docker_stop(remove_volumes=True)