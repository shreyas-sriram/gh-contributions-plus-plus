[![Actions Status](https://github.com/shreyas-sriram/gh-contributions-plus-plus/workflows/CI/badge.svg)](https://github.com/shreyas-sriram/gh-contributions-plus-plus/actions)
# gh-contributions-plus-plus

A Golang application to aggregate the GitHub contributions of multiple accounts.

![gh-contributions](https://npdt6zg2ga.execute-api.ap-south-1.amazonaws.com/dev/aggregate?username=shreyas-sriram&username=shreyas-mesh7&year=2020&theme=light)

### GitHub Profile README Integration

This application is currently deployed and can be used as follows in your GitHub profile README<br>
```
![gh-contributions](https://npdt6zg2ga.execute-api.ap-south-1.amazonaws.com/dev/aggregate?username={username1}&username={username2}&year=2020&theme=dark)
```

### Usage

1. Start the server

- on host machine
```
make build.x && make server.start
```
> where 'x' is the runtime environment. Currently, `linux` and `mac` are available in the Make target.

- in docker
```
make docker
```

2. Open browser and browse to
```
localhost:3000/api/contributions?username={username}&username={username}&year={year}&theme={theme}
```

#### Options

```
username
  description : list of usernames

year
  description : year
  default     : current year

theme
  description : theme
  values      : <light | dark>
  default     : light
```
