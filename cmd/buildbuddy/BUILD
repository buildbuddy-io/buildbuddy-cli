load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["buildbuddy.go"],
    importpath = "github.com/buildbuddy-io/buildbuddy-cli/cmd/buildbuddy",
    visibility = ["//visibility:private"],
    deps = [
        "//sidecar:go_default_library",
        "@com_github_bazelbuild_bazelisk//core:go_default_library",
        "@com_github_bazelbuild_bazelisk//repositories:go_default_library",
    ],
)

go_binary(
    name = "buildbuddy",
    embed = [":go_default_library"],
    visibility = ["//visibility:public"],
)
