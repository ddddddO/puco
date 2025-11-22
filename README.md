# 💊PuCo
Tool that aids ***P***HP***U***nit and P***CO***V

Running tests takes a lot of time, and generating test coverage also takes a significant amount of time.
This tool might be useful in such situations.

![](./assets/puco_demo.gif)

> [!IMPORTANT] 
> `vendor/bin/phpunit` and `pcov` are required.
> e.g. `sudo apt-get update && sudo apt-get install phpX.Y-pcov` (X.Y is the PHP version)

## With this tool,
1. Select test files to run,
1. Select files for which you want to generate coverage reports (HTML),
1. You might be able to execute steps 1 and 2 easily and quickly. Probably, probably..

## Which test files can be selected?
Files located in `tests/` directory.

## Which files can be selected to get coverage?
Files located in `src/` or `app/` directory.

## How to install?

### Homebrew

```console
brew install ddddddO/tap/puco
```

### Go
```console
go install github.com/ddddddO/puco/cmd/puco@latest
```

## Usage

```console
$ puco --help
Usage: puco [options]
puco

Options:
  -repeat
        This flag starts with data selected by the most recently executed puco.

Example:
  puco          # normal launch
  puco --repeat # launch using the most recent data

Processing description:
  1. You can select multiple test files to run (fuzzy search available).
  2. You can select multiple PHP files for which you want to calculate coverage (fuzzy search available).
  3. Calculate the longest matching directory path from multiple selected PHP file paths in step 2
    - ※ Note that only the PHP file paths selected in step 2 are not the target for coverage calculation. Instead, the directory path under the longest match calculated becomes the target for coverage calculation. If there are numerous PHP files under the calculated directory path, the coverage calculation process may become slow.
  4. If Steps 1 and 3 and an existing phpunit.xml are present, generate phpunitxml_generated_by_puco.xml based on them.
  5. Assemble and execute the php command.
  6. Coverage reports are generated under the coverage-puco directory.

WARNING:
  When puco is run for the first time, a configuration file named ~/.config/puco.toml is created. This configuration file contains a key: CommandToSpecifyBeforePHPCommand. It specifies that the PHP command should be executed via the Docker command. If you wish to execute the PHP command directly, please set the value of this key to "" or delete this entire line.
$
```

## Processing of PuCo (en)

1. You can select multiple test files to run (fuzzy search available).
1. You can select multiple PHP files for which you want to calculate coverage (fuzzy search available).
1. Calculate the longest matching directory path from multiple selected PHP file paths in step 2
    - ※ Note that only the PHP file paths selected in step 2 are not the target for coverage calculation. Instead, the directory path under the longest match calculated becomes the target for coverage calculation. **If there are numerous PHP files under the calculated directory path, the coverage calculation process may become slow.**
1. If Steps 1 and 3 and an existing `phpunit.xml` are present, generate `phpunitxml_generated_by_puco.xml` based on them.
1. Assemble and execute the `php` command.
1. Coverage reports are generated under the `coverage-puco` directory.

> [!WARNING]
> ※ When `puco` is run for the first time, a configuration file named `~/.config/puco.toml` is created.
> This configuration file contains a key: `CommandToSpecifyBeforePHPCommand`. It specifies that the PHP command should be executed via the Docker command. If you wish to execute the PHP command directly, please set the value of this key to `""` or delete this entire line.

<details><summary>Processing of PuCo (ja)</summary>

1. 実行したいテストファイルを複数選択できます（fuzzyに検索可能）
1. カバレッジを算出したいPHPファイルを複数選択できます（fuzzyに検索可能）
1. 2で複数選択されたPHPファイルパスから最長一致のディレクトリパスを計算
    - ※ 2で選択された各PHPファイルパスのみがカバレッジ計算対象では無く、算出された最長一致のディレクトリパス配下がカバレッジ算出対象になることに注意してください。**算出されるディレクトリ配下に多数のPHPファイルがある場合、カバレッジ計算処理が遅くなるかもしれません。**
1. 1と3と既存の`phpunit.xml`があれば、それらを元に`phpunitxml_generated_by_puco.xml`を生成
1. 実行する`php`コマンドを組み立て、実行する
1. `coverage-puco`ディレクトリ配下にカバレッジレポートが生成される

> [!WARNING]
> ※`puco`初回実行時に、`~/.config/puco.toml`という設定ファイルができます。
> この設定ファイル内のキー:`CommandToSpecifyBeforePHPCommand`にdockerコマンド越しにphpコマンドを実行するよう記載していますが、直接phpコマンドを実行したい場合は、このキーの値を`""`にしていただくか、この行ごと消してください。

</details>

## TODO
- [ ] カバレッジレポートをHTML形式以外でも出力できるようにする
- [ ] ヒストリー機能欲しい
    - 何度もファイル選択は手間。ただ、ツールで組み立てられたコマンドは表示されるので、それコピペで実行でも代替できるから後でいいかも
    - 一旦repeatフラグを実装したから、ほんとに欲しくなってからでいいかも