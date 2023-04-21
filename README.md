# Organize Repo Folder

This utility is written in Go and features a simple script to organize your local repositories according to this opinionated structure (brutally):

```shell
Repos
├── bjerkio
│   ├── jobs
│   ├── oidc-react
│   ├── pulumi-appsync-modules
│   └── vault-action
├── cobraz
│   ├── firestore-graphql-scalars
│   ├── norwegian-ssn-graphql-scalar
│   ├── organize-repo-folder
│   └── pulumi-action
└── sindresorhus
    ├── awesome
    ├── awesome-nodejs
    └── gifski
```

In the example above, `bjerkio`, `cobraz` and `sindresorhus` are organizations or GitHub users. The scripts get this from the git source (ex. `git@github.com:*cobraz*/organize-repo-folder.git`). Please note, we are lowercasing both the organization name and repository name. Repositories that either has no origin source have a non-standard GitHub path or maybe having issues are parsed by Git, will be kept in a backup directory.

## Warning

This command-line utility is very bare-bone. You should make sure that you have a backup of everything. Remember that probably all your repositories will get a new location, so – that
could cause problems for running applications and editors. You probably want to close them before running any of this.

Also, as we are creating a `(name of your repository-folder)-bak` of your old folder and consequently moving all your repositories over to a new folder, you could risk losing stuff if the `-bak` folder is there before running.

Lastly, we do not confirm anything. Once you run the command – here be dragons!

## Install

```shell
brew install cobraz/tools/organize-repo-folder
```

**Notes**: The library is not tested on Linux or Windows. There are [executables available](https://github.com/cobraz/organize-repo-folder/releases/latest) at every release >1.0.3

## Help

```shell
> organize-repo-folder --help
NAME:
   organize-repo-folder - moves every repo into each github org

USAGE:
   main [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --path value  path to your repositories (default: "<your home directory>/Repos")
   --dry-run     (default: false)
   --help, -h    show help (default: false)
```

## Other projects

I personally love [ghq](https://github.com/x-motemen/ghq). I recommend it!

## Contribute

Please, oh pretty please do contribute! If you feel this helps you out, but you want to increase the quality of this software, please submit pull requests. Look at our issues page for more information – as previously stated, it's pretty bare-bone. Making it faster, better is something everyone wants.
