load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["autoconfig.go"],
    importpath = "github.com/buildbuddy-io/buildbuddy-cli/autoconfig",
    visibility = ["//visibility:public"],
    deps = [
        "//commandline:go_default_library",
    ],
)
