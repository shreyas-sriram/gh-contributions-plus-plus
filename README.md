[![Actions Status](https://github.com/shreyas-sriram/gh-contributions-aggregator/workflows/CI/badge.svg)](https://github.com/shreyas-sriram/gh-contributions-aggregator/actions)
# gh-contributions-aggregator

A simple application to aggregate the GitHub contributions of multiple accounts.

### Usage

Start the server
```
make build.x && make server.start
```
> where 'x' is the runtime environment. Currently, `linux` and `mac` are available in the Make target.

Open browser and browse to
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
