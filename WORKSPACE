workspace(
    name = "buildbuddy_cli",
    managed_directories = {"@npm": ["node_modules"]},
)

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

# Go

http_archive(
    name = "io_bazel_rules_go",
    sha256 = "207fad3e6689135c5d8713e5a17ba9d1290238f47b9ba545b63d9303406209c6",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.24.7/rules_go-v0.24.7.tar.gz",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.24.7/rules_go-v0.24.7.tar.gz",
    ],
)

http_archive(
    name = "bazel_gazelle",
    sha256 = "86c6d481b3f7aedc1d60c1c211c6f76da282ae197c3b3160f54bd3a8f847896f",
    urls = [
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.19.1/bazel-gazelle-v0.19.1.tar.gz",
    ],
)

load("@io_bazel_rules_go//go:deps.bzl", "go_download_sdk", "go_register_toolchains", "go_rules_dependencies")

go_download_sdk(
    name = "go_sdk",
    version = "1.15",
)

go_rules_dependencies()

go_register_toolchains()

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies", "go_repository")

gazelle_dependencies()

# Node

http_archive(
    name = "build_bazel_rules_nodejs",
    sha256 = "452bef42c4b2fbe0f509a2699ffeb3ae2c914087736b16314dbd356f3641d7e5",
    urls = ["https://github.com/bazelbuild/rules_nodejs/releases/download/2.3.0/rules_nodejs-2.3.0.tar.gz"],
)

load("@build_bazel_rules_nodejs//:index.bzl", "yarn_install")

# Docker

http_archive(
    name = "io_bazel_rules_docker",
    sha256 = "dc97fccceacd4c6be14e800b2a00693d5e8d07f69ee187babfd04a80a9f8e250",
    strip_prefix = "rules_docker-0.14.1",
    urls = ["https://github.com/bazelbuild/rules_docker/releases/download/v0.14.1/rules_docker-v0.14.1.tar.gz"],
)

load(
    "@io_bazel_rules_docker//toolchains/docker:toolchain.bzl",
    docker_toolchain_configure = "toolchain_configure",
)

docker_toolchain_configure(
    name = "docker_config",
    docker_path = "/usr/bin/docker",
)

load(
    "@io_bazel_rules_docker//repositories:repositories.bzl",
    container_repositories = "repositories",
)

container_repositories()

load(
    "@io_bazel_rules_docker//go:image.bzl",
    _go_image_repos = "repositories",
)

_go_image_repos()

# Proto

# NB: The name must be "com_google_protobuf".
http_archive(
    name = "com_google_protobuf",
    sha256 = "1672819a0baf3c57e2ab96bc7cd9935f8b58c0172317c44aa44722d4b1b30f8b",
    strip_prefix = "protobuf-3.11.2",
    urls = ["https://github.com/protocolbuffers/protobuf/releases/download/v3.11.2/protobuf-all-3.11.2.zip"],
)

load("@com_google_protobuf//:protobuf_deps.bzl", "protobuf_deps")

protobuf_deps()

# Pick your poison:
#  - either use the git_repository below to pull Buildbuddy from Github
#  - or uncomment local_repository to use a buildbuddy directory parallel
#    to your buildbuddy-internal directory. This makes cross repo
#    development easier.

load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository")

git_repository(
    name = "com_github_buildbuddy_io_buildbuddy",
    commit = "feeff7d86fc266c93f05c46f7dab5d9462cca7f7",  # autoupdate buildbuddy-io/buildbuddy
    remote = "https://github.com/buildbuddy-io/buildbuddy.git",
)

# local_repository(
#     name = "com_github_buildbuddy_io_buildbuddy",
#     path = __workspace_dir__ + "/../buildbuddy",
# )

load("@com_github_buildbuddy_io_buildbuddy//:deps.bzl", "install_buildbuddy_dependencies")

install_buildbuddy_dependencies()

yarn_install(
    name = "npm",
    package_json = "@com_github_buildbuddy_io_buildbuddy//:package.json",
    yarn_lock = "@com_github_buildbuddy_io_buildbuddy//:yarn.lock",
)

# @bazel/labs (for ts_proto_library)

load("@npm//@bazel/labs:package.bzl", "npm_bazel_labs_dependencies")

npm_bazel_labs_dependencies()

# We need to explicitly import this go_repo in buildbuddy-internal because Google cloud go
# repositories share an internal dependency, and will not build without this.
# See https://github.com/bazelbuild/bazel-gazelle/issues/619 for more.

go_repository(
    name = "com_google_cloud_go_storage",
    importpath = "cloud.google.com/go/storage",
    sum = "h1:RPUcBvDeYgQFMfQu1eBMq6piD1SXmLH+vK3qjewZPus=",
    version = "v1.5.0",
)

# BuildBuddy Toolchain

http_archive(
    name = "io_buildbuddy_buildbuddy_toolchain",
    sha256 = "9055a3e6f45773cd61931eba7b7cf35d6477ab6ad8fb2f18bf9815271fc682fe",
    strip_prefix = "buildbuddy-toolchain-52aa5d2cc6c9ba7ee4063de35987be7d1b75f8e2",
    urls = ["https://github.com/buildbuddy-io/buildbuddy-toolchain/archive/52aa5d2cc6c9ba7ee4063de35987be7d1b75f8e2.tar.gz"],
)

load("@io_buildbuddy_buildbuddy_toolchain//:deps.bzl", "buildbuddy_deps")

buildbuddy_deps()

load("@io_buildbuddy_buildbuddy_toolchain//:rules.bzl", "buildbuddy")

buildbuddy(name = "buildbuddy_toolchain")

# gazelle:repository go_repository name=org_golang_x_sys importpath=golang.org/x/sys
