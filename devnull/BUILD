load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["devnull.go"],
    importpath = "github.com/buildbuddy-io/buildbuddy-cli/devnull",
    visibility = ["//visibility:public"],
    deps = [
        "@com_github_buildbuddy_io_buildbuddy//proto:publish_build_event_go_proto",
        "@com_github_buildbuddy_io_buildbuddy//server/interfaces:go_default_library",
    ],
)
