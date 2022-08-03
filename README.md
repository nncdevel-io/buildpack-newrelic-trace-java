# paketo-newrelic-trace-java

[New Relic Java Agent](https://docs.newrelic.com/jp/docs/apm/agents/java-agent/installation/install-java-agent/)の利用に必要な `newrelic-java-agent.jar` を挿入するビルドパックです。

## ゴール

Newrelic Java Agent では Jar ファイルが提供されますが、バイトコード変更を行う JavaAgent の Jar として提供されているため、java コマンドの`-jar` オプションではなく、`-javaagent`オプションに渡す必要があります。

本プロジェクトのアプリケーションは SpringBoot 2.3 からサポートされた CloudNativeBuildpacks を利用した OCI イメージビルドを利用しているため、このイメージビルド時に jar ファイルを挿入し、java コマンド起動オプションを設定する必要があります。

## 設定

| 環境編集                     | 説明                                                                                                                   |
| ---------------------------- | ---------------------------------------------------------------------------------------------------------------------- |
| `$BP_USE_NEWRELIC`           | New Relic Java Agenct Jarを挿入するスイッチです。 デフォルトでは `false` が設定され、利用されません。 |
| `$BP_NEWRELIC_AGENT_VERSION` | 挿入する New Relic Java Agent Jar のバージョンを指定します。 <br> デフォルトでは `7.9.0`(22/08/02 現在) を利用します。 |

## ビルドに必要なもの

-   docker

golang、pack CLI をローカルにインストールする場合は下記の 2 つが必要になる。

-   [pack CLI](https://buildpacks.io/docs/install-pack/)
-   [golang](https://golang.org/doc/install)

インストールしない場合は docker-compose をインストールしてください。

## docker-compose を利用したビルド手順

まずは `docker-compose.yml.example` を `docker-compose.yml` としてコピーします。

プロキシ設定が必要な環境ではコピーしたファイルを開き、`http_proxy`、`https_proxy` の認証情報を各人のもので上書き、保存してください。

golang のビルド

```bash
$ docker-compose run build-golang
```

buildpack のパッケージング(Docker イメージ化)

```bash
$ docker-compose run package-buildpack-image
```

イメージの push

```bash
$ docker tag paketo-newrelic-java-agent nncdevel-io/paketo-newrelic-java-agent:0.0.1
$ docker push nncdevel-io/paketo-newrelic-java-agent:0.0.1
```

## golang、pack-cli をインストールしている場合のビルド手順

golang のビルド

```bash
$ ./scripts/build-golang.sh
```

buildpack のパッケージング(Docker イメージ化)

```bash
$ ./scripts/package-buildpack-image.sh
```

イメージの push

```bash
$ ./scripts/push-image.sh
```
