import sys
import os
import json

config_file = os.path.join(os.path.expanduser('~'), ".ssh", "accountswitcherconfig.json")
temp_config = os.path.join(os.path.expanduser('~'), ".ssh", "accountswitcherconfig.txt")

def read_config_file():
    with open(config_file) as f:
        config = json.load(f)
    return config

def write(account: str):
    with open(temp_config, 'w') as f:
        f.write(account.strip())


def fix_url(account: str, repo: str):
    data = conf[account]
    prefix = data['prefix']

    owner, name = repo.split("/")
    url = f'{prefix}:{owner}/{name}.git'
    print(url)

    return url 

def change_account(account: str):
    # write account to temp_config
    write(account)

    email_command = "git config --global user.email {}"
    data = conf[account]
    email = data['email']
    command = email_command.format(email)
    print(command)
    run_command(command)

def clone_repo(url):
    command = f"git clone {url}"
    run_command(command)

def run_command(command: str):
    print('running command => ', command)
    os.system(command)

def info():
    with open(temp_config) as f:
        print(f.read())

if __name__ == "__main__":
    # read from sys
    conf = read_config_file()
    if len(sys.argv) > 2:
        account = sys.argv[1] # work/personal
        repo = sys.argv[2] # ybgirgin3/repo-url-account-dedicator
        change_account(account=account)
        url = fix_url(account, repo)
        clone_repo(url)
    else:
        if sys.argv[1] in ('personal', 'work'):
            account = sys.argv[1] # work/personal
            change_account(account=account)
        elif sys.argv[1] in ('info'):
            info()


