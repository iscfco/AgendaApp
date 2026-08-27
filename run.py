import sys
import subprocess

from artifacts.run.build import execute_docker_build
from artifacts.run.run import execute_docker_run
from artifacts.run.metrics import execute_metrics

# List of available commands
COMMANDS = {
    "build": execute_docker_build,
    "run": execute_docker_run,
    "stop": 'docker-compose -f .\\artifacts\\docker-compose.yaml down',
    "lint": '...',
    "test": '...',
    "metrics": execute_metrics,
}

def print_help():
    print("\nAvailable commands:")
    for cmd in COMMANDS:
        print(f"- python run.py {cmd}")
    print()

def main():
    # If no command is provided print help and exit
    if len(sys.argv) < 2:
        print_help()
        sys.exit(1)

    # Get commnand and process
    action = sys.argv[1]
    
    # If the command is not supported, print error, show help and exit
    if action not in COMMANDS:
        print(f"❌ Error: Command >>> '{action}' <<< not supported.")
        print_help()
        sys.exit(1)
        
    # Run the command
    function = COMMANDS[action]
    print(f"🚀 Running: >>> {action} <<<:\n")
    
    try:
        function()
    except subprocess.CalledProcessError:
        print(f"\n❌ Execution failed for command >>> '{action}' <<<\n")
        sys.exit(1)

if __name__ == "__main__":
    main()
