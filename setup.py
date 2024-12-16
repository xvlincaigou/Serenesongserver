import os
import subprocess
import shutil
import sys
from git import Repo

def CloneRepo(github_url, local_dir, target_dir, new_name):
    try:
        # Ensure the local directory does not already exists
        if os.path.exists(os.path.join(local_dir, new_name)):
            print(f"Local directory {os.path.join(local_dir, new_name)} already exists. Skipping clone.")
            return
        # Create the local directory
        os.makedirs(local_dir)
        # Initialize a new repository in the local directory
        print(f"Initializing sparse checkout for {github_url}...")
        repository = Repo.init(local_dir)
        remote_org = repository.create_remote('origin', github_url)
        # Fetch remote branches to determine the default branch
        remote_org.fetch()
        branches = [ref.name.split('/')[-1] for ref in remote_org.refs]
        print(f"Available branches: {branches}")
        # Determine the branch to use (prioritize master, then main)
        checkout = 'master' if 'master' in branches else 'main'
        if checkout not in branches:
            raise Exception("Neither 'master' nor 'main' branch exists in the repository.")
        print(f"Using branch: {checkout}")
        # Enable sparse checkout
        sparse_path = os.path.join(local_dir, '.git', 'info', 'sparse-checkout')
        with open(sparse_path, 'w') as f:
            f.write(target_dir + '\n')
        print(f"Configured sparse checkout for folder: {target_dir}.")
        # Fetch and check out the specific folder
        remote_org.fetch()
        repository.git.config('core.sparseCheckout', 'true')
        repository.git.checkout(checkout)
        remote_org.pull(checkout)
        print(f"Cloned folder '{target_dir}' from branch '{checkout}' into {local_dir}.")
        # Remove the .git folder
        git_folder = os.path.join(local_dir, ".git")
        if os.path.exists(git_folder):
            shutil.rmtree(git_folder)
        else:
            print(f"[!].git folder not found in {repo_path}.")
        # Rename the folder to match the target name
        old_path = os.path.join(local_dir, target_dir)
        new_path = os.path.join(local_dir, new_name)
        if os.path.exists(old_path):
            os.rename(old_path, new_path)
            print(f"Renamed folder '{old_path}' to '{new_path}'")
        else:
            raise FileNotFoundError(f"The target folder '{old_path}' does not exist after clone.")
    
    except Exception as e:
        print(f"Error cloning repository: {e}")
        sys.exit(1)
        
def RunScripts(scripts):
    for script in scripts:
        print(f"Running {script}...")
        try:
            result = subprocess.run(
                ["python", script], 
                check=True,
                capture_output=True, 
                text=True
            )
            print(f"Output of {script}: {result.stdout}")
        except subprocess.CalledProcessError as e:
            print(f"Error while running {script}: {e.stderr}")
            break

if __name__ == "__main__":
    data_songci = "https://github.com/chinese-poetry/chinese-poetry.git"
    data_rhymes = "https://github.com/charlesix59/chinese_word_rhyme.git"
    db_dir = "./database"
    # Download and preprocess the data
    CloneRepo(data_songci, db_dir, "宋词", "songci")
    CloneRepo(data_rhymes, db_dir, "data", "rhymes")
    # Run the scripts to generate the database
    scripts = [
        "./scripts/rev_yunshu.py",
        "./scripts/load_authors.py",
        "./scripts/load_ci.py",
        "./scripts/load_cipai.py",
        "./scripts/load_pingze.py",
        "./scripts/load_yunshu.py",
        "./scripts/load_yunshu_rev.py",
    ]
    RunScripts(scripts)