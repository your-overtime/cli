# Your Overtime CLI (otcli)

You can download the newest build [via](https://github.com/your-overtime/cli/releases/latest) and move the file to `/usr/local/bin/otcli` or use for linux:

```bash
curl https://api.github.com/repos/your-overtime/cli/releases/latest | grep otcli_linux | grep -v 'arm\|name' | awk -F'"' '{print$4}' | wget -i - && sudo mv otcli_linux /usr/local/bin/otcli && chmod ugo+rwx /usr/local/bin/otcli
```

The `otcli` needs the following configuration file, which can be generated with `otcli conf init`.

```json
{
 "Host": "https://your-overtime.de",
 "Token": "token secretGeneratedToken",
 "DefaultActivityDesc": "Coding"
}

```
