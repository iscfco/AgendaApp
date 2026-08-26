import subprocess

def execute_docker_run():
    """Run all tasks required to run docker compose services"""

    print("📋 [Build] Running services")

    # Docker Compose Command
    docker_cmd = 'docker-compose -f .\\artifacts\\docker-compose.yaml up'
    
    #check=True will raise an error if the command fails
    subprocess.run(docker_cmd, shell=True, check=True)

