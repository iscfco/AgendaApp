import subprocess

def execute_docker_build():
    """Run all tasks required to build the images for the project"""

    print("📋 [Build] Running images build")

    # Docker Compose Command
    docker_cmd = 'docker-compose -f .\\artifacts\\docker-compose.yaml build --no-cache'
    
    #check=True will raise an error if the command fails
    subprocess.run(docker_cmd, shell=True, check=True)
    
    print("✨ [Build] Images built successfully")
