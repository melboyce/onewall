# onewall

> just the one, thanks

`onewall` reads an existing `~/.fehbg` file, parses the wallpaper list, and
replaces one with a randomly chosen wallpaper :sparkles:

If you run Linux multihead and set your wallpaper with `feh`, this could be
useful. Although a shell script will do just as well :woman_shrugging:

# usage

```bash
# usage:
# onewall [-pos N] [-l|-p] <dir> [<dir> [<dir>...]]
#   -pos N  try to set a wall for a specific xinerama screen
#   -l      restrict to landscape
#   -p      restrict to portrait

# change the wall on head 2
onewall -pos 1 /some/dir

# change the wall on head 1, only landscape
onewall -pos 0 -l /a/bunch/of/walls /another/dir
```
