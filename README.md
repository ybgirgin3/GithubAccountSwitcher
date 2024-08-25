# GithubAccountSwitcher

config_file = os.path.join(os.path.expanduser('~'), ".ssh", "accountswitcherconfig.json")
temp_config = os.path.join(os.path.expanduser('~'), ".ssh", "accountswitcherconfig.txt")

## REQUIRED FILES
`~/.ssh/accountswitcherconfig.json`
```json
{
  "personal": { "email": "personal@email.com", "name": "Yusuf Berkay Girgin", "prefix": "git@personal.github.com" },
  "work": {
    "email": "work@email.com",
    "name": "Yusuf Berkay Girgin",
    "prefix": "git@work.github.com"
  }
}

```


