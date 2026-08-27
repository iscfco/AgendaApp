import subprocess

def execute_metrics():
    """Showing metrics of the repo
    Currently on works using git bash, we need to run it through a linux container"""

    print("📋 [Metrics] Running metrics")


    # Comando corregido usando ForEach-Object
    cmd = 'git ls-files | ForEach-Object { Get-Content $_ } | Measure-Object -Line'

    try:
        subprocess.run(["powershell", "-Command", cmd], check=True)
    except subprocess.CalledProcessError as e:
        print(f"❌ Error al ejecutar las métricas: {e}")

