# Todo-Mitsukeru-Kun
Compile a list of TODO comments for each file into an issue.

## Usage
Create a `.github/workflows/todo_mitsukeru_kun.yaml` like this.
```yaml
on:
  schedule:
    - cron:  '0 9 * * 1'

jobs:
  todo_mitsukeru_kun:
    runs-on: ubuntu-latest
    name: test
    steps:
      - name: Checkout
        uses: actions/checkout@v3
      - name: Use todo-mitsukeru-kun
        uses: GOD-oda/todo-mitsukeru-kun@v0.0.1
        env:
          INPUT_GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
          INPUT_TARGET_DIR: "src"
```

This workflow runs on Monday morning at 9:00 a.m. and searches for TODO comments in the `INPUT_TARGET_DIR` directory.


## Contribution

Please run `cd src && sh ./build.sh` and include the created binary in the PR.