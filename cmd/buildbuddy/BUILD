load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["buildbuddy.go"],
    importpath = "github.com/buildbuddy-io/buildbuddy-cli/cmd/buildbuddy",
    visibility = ["//visibility:private"],
    deps = [
        "//:go_default_library",
        "//commandline:go_default_library",
        "//logging:go_default_library",
        "//parser:go_default_library",
        "//sidecar:go_default_library",
        "@com_github_bazelbuild_bazelisk//core:go_default_library",
        "@com_github_bazelbuild_bazelisk//repositories:go_default_library",
        "@com_github_buildbuddy_io_buildbuddy//server/version:go_default_library",
    ],
)

go_binary(
    name = "buildbuddy",
    embed = [":go_default_library"],
    visibility = ["//visibility:public"],
)

go_binary(
    name = "buildbuddy-darwin",
    out = "buildbuddy-darwin-amd64",
    embed = [":go_default_library"],
    gc_linkopts = [
        "-s",
        "-w",
    ],
    goarch = "amd64",
    goos = "darwin",
    pure = "on",
    tags = ["manual"],
    visibility = ["//visibility:public"],
)

go_binary(
    name = "buildbuddy-linux",
    out = "buildbuddy-linux-amd64",
    embed = [":go_default_library"],
    gc_linkopts = [
        "-s",
        "-w",
    ],
    goarch = "amd64",
    goos = "linux",
    pure = "on",
    tags = ["manual"],
    visibility = ["//visibility:public"],
)

go_binary(
    name = "buildbuddy-linux-arm64",
    out = "buildbuddy-linux-arm64",
    embed = [":go_default_library"],
    gc_linkopts = [
        "-s",
        "-w",
    ],
    goarch = "arm64",
    goos = "linux",
    pure = "on",
    tags = ["manual"],
    visibility = ["//visibility:public"],
)

go_binary(
    name = "buildbuddy-windows",
    out = "buildbuddy-windows-amd64.exe",
    embed = [":go_default_library"],
    goarch = "amd64",
    goos = "windows",
    pure = "on",
    tags = ["manual"],
    visibility = ["//visibility:public"],
)
